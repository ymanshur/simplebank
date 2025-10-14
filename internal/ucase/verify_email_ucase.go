package ucase

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/validator"
	"github.com/ymanshur/simplebank/pkg/util"
)

type VerifyEmailUcase interface {
	Create(ctx context.Context, req CreateVerifyEmailRequest) (*CreateVerifyEmailResponse, error)
}

type verifyEmailUcase struct {
	repo repo.Repo
}

func NewVerifyEmailUseCase(repo repo.Repo) VerifyEmailUcase {
	return &verifyEmailUcase{repo: repo}
}

type CreateVerifyEmailRequest struct {
	Username string `json:"username"`
}

func (r CreateVerifyEmailRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.By(validator.ValidUsername)),
	)
}

type CreateVerifyEmailResponse struct {
	ID         int64
	Username   string
	Email      string
	FullName   string
	SecretCode string
	IsUsed     bool
	CreatedAt  time.Time
	ExpiredAt  time.Time
}

func (u *verifyEmailUcase) Create(ctx context.Context, req CreateVerifyEmailRequest) (*CreateVerifyEmailResponse, error) {
	user, err := u.repo.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			err = typex.ErrDataNotFound("user")
			// Make compensation if DB transaction need more than
			// the time processor takes the task
			// err = util.JoinErrors(err, asynq.SkipRetry)

			return nil, err
		}

		return nil, errors.Wrap(err, "failed to get user")
	}

	verifyEmail, err := u.repo.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create verify email")
	}

	rsp := convertVerifyEmail(verifyEmail, user)
	return &rsp, nil
}

func convertVerifyEmail(verifyEmail db.VerifyEmail, user db.User) CreateVerifyEmailResponse {
	return CreateVerifyEmailResponse{
		ID:         verifyEmail.ID,
		Username:   verifyEmail.Username,
		Email:      verifyEmail.Email,
		FullName:   user.FullName,
		SecretCode: verifyEmail.SecretCode,
		IsUsed:     verifyEmail.IsUsed,
		CreatedAt:  verifyEmail.CreatedAt,
		ExpiredAt:  verifyEmail.ExpiredAt,
	}
}
