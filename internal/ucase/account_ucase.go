package ucase

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/validator"
)

type accountUcase struct {
	repo repo.Repo
}

func NewAccountUseCase(repo repo.Repo) AccountUseCase {
	return &accountUcase{
		repo: repo,
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (r AuthRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.By(validator.ValidUsername)),
		validation.Field(&r.Role, validation.In(common.BankerRole, common.DepositorRole)),
	)
}

type CreateAccountRequest struct {
	Auth     AuthRequest `json:"auth"`
	Currency string      `json:"currency"`
}

func (r CreateAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Auth),
		validation.Field(&r.Currency, validation.Required, validation.By(validator.ValidCurrency)),
	)
}

type AccountResponse struct {
	ID        int64
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}

func (u *accountUcase) Create(ctx context.Context, req CreateAccountRequest) (*AccountResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	if req.Auth.Username == "" {
		return nil, typex.ErrUnAuthorized("unknown user")
	}

	// TODO: Create account with default balance,
	// since there is not top-up API
	arg := db.CreateAccountParams{
		Owner:    req.Auth.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := u.repo.CreateAccount(ctx, arg)
	if err != nil {
		errCode := repo.ErrorCode(err)
		if errCode == repo.ForeignKeyViolation || errCode == repo.UniqueViolation {
			return nil, typex.ErrUnProcessableEnity(err.Error())
		}

		return nil, errors.Wrap(err, "create account")
	}

	rsp := AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
	return &rsp, nil
}

type GetAccountRequest struct {
	Auth AuthRequest `json:"auth"`
	ID   int64       `json:"id"`
}

func (r GetAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Auth),
		validation.Field(&r.ID, validation.Required, validator.ValidID()),
	)
}

func (u *accountUcase) Get(ctx context.Context, req GetAccountRequest) (*AccountResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	if req.Auth.Username == "" {
		return nil, typex.ErrUnAuthorized("unknown user")
	}

	account, err := u.repo.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			return nil, typex.NewErrDataNotFound("account")
		}

		return nil, errors.Wrap(err, "get account")
	}

	if account.Owner != req.Auth.Username {
		return nil, typex.ErrForbidden("account doesn't belong to the authorized user")
	}

	rsp := AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
	return &rsp, nil
}

type PagingRequest struct {
	ID   int32 `json:"id"`
	Size int32 `json:"size"`
}

func (r PagingRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, validator.ValidID()),
		validation.Field(&r.Size, validation.Required, validation.By(validator.ValidPageSize)),
	)
}

type ListAccountRequest struct {
	Auth AuthRequest   `json:"auth"`
	Page PagingRequest `json:"page"`
}

func (r ListAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Auth),
		validation.Field(&r.Page),
	)
}

func (u *accountUcase) List(ctx context.Context, req ListAccountRequest) ([]AccountResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	if req.Auth.Username == "" {
		return nil, typex.ErrUnAuthorized("unknown user")
	}

	arg := db.ListAccountsParams{
		Owner: pgtype.Text{
			String: req.Auth.Username,
			Valid:  true,
		},
		Limit:  req.Page.Size,
		Offset: (req.Page.ID - 1) * req.Page.Size,
	}

	accounts, err := u.repo.ListAccounts(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "list accounts")
	}

	rsp := convertAccounts(accounts)
	return rsp, nil
}

func convertAccount(account db.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func convertAccounts(accounts []db.Account) []AccountResponse {
	var res []AccountResponse
	for _, account := range accounts {
		res = append(res, convertAccount(account))
	}
	return res
}
