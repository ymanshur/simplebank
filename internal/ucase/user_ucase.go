package ucase

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/validator"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserUniqueViolation = typex.ErrUnProcessableEnity("user unique constraint violated")
	ErrPermisionDenied     = typex.ErrForbidden("permission denied")
)

type userUcase struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewUserUseCase(
	config util.Config,
	store db.Store,
	tokenMaker token.Maker,
	taskDistributor worker.TaskDistributor,
) UserUseCase {
	return &userUcase{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.By(validator.ValidUsername)),
		validation.Field(&r.FullName, validation.Required, validation.By(validator.ValidFullName)),
		validation.Field(&r.Email, validation.Required, validation.By(validator.ValidEmail)),
		validation.Field(&r.Password, validation.Required, validation.By(validator.ValidPassword)),
	)
}

type UserResponse struct {
	Username          string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

func (u *userUcase) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	if err := validation.Validate(&req); err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password")
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			FullName:       req.FullName,
			Email:          req.Email,
			HashedPassword: hashedPassword,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := worker.PayloadSendVerifyEmail{Username: user.Username}
			taskOpts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return u.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &taskPayload, taskOpts...)
		},
	}

	txResult, err := u.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, typex.ErrUnProcessableEnity("user unique constraint violated")
		}

		return nil, errors.Wrap(err, "create user")
	}

	rsp := convertUser(txResult.User)
	return &rsp, nil
}

type LoginUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	ClientIp  string `json:"client_ip"`
}

func (r LoginUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.By(validator.ValidUsername)),
		validation.Field(&r.Password, validation.Required, validation.By(validator.ValidPassword)),
	)
}

type LoginUserResponse struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  UserResponse
}

func (u *userUcase) Login(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error) {
	if err := validation.Validate(&req); err != nil {
		return nil, err
	}

	user, err := u.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, typex.NewErrDataNotFound("user")
		}

		return nil, errors.Wrap(err, "get user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, typex.ErrUnAuthorized("incorrect password")
		}

		return nil, errors.Wrap(err, "check password")
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		u.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create access token")
	}

	refreshToken, refreshPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		u.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create refresh token")
	}

	session, err := u.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    req.UserAgent,
		ClientIp:     req.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create session")
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  convertUser(user),
	}
	return &rsp, nil
}

type UpdateUserRequest struct {
	Auth     AuthRequest `json:"auth"`
	Username string      `json:"username"`
	FullName string      `json:"full_name"`
	Email    string      `json:"email"`
	Password string      `json:"passoword"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Auth),
		validation.Field(&r.Username, validation.By(validator.ValidUsername)),
		validation.Field(&r.FullName, validation.By(validator.ValidFullName)),
		validation.Field(&r.Email, validation.By(validator.ValidEmail)),
		validation.Field(&r.Password, validation.By(validator.ValidPassword)),
	)
}

func (u *userUcase) Update(ctx context.Context, req UpdateUserRequest) (*UserResponse, error) {
	if !common.HasPermission(req.Auth.Role, []string{util.BankerRole, util.DepositorRole}) {
		return nil, ErrPermisionDenied
	}

	if err := validation.Validate(&req); err != nil {
		return nil, err
	}

	if req.Auth.Username == "" {
		return nil, typex.ErrUnAuthorized("unknown user")
	}

	if req.Auth.Role != util.BankerRole && req.Auth.Username != req.Username {
		return nil, typex.ErrForbidden("cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: pgtype.Text{
			String: req.FullName,
			Valid:  req.FullName != "",
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Username)
		if err != nil {
			return nil, errors.Wrap(err, "hash password")
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := u.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, typex.NewErrDataNotFound("user")
		}
		return nil, errors.Wrap(err, "update user")
	}

	rsp := convertUser(user)
	return &rsp, nil
}

func convertUser(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
