package ucase

import (
	"context"

	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
)

type UseCase struct {
	User        UserUseCase
	Token       TokenUseCase
	Account     AccountUseCase
	Transaction TransactionUseCase
}

func NewUseCase(
	config util.Config,
	store db.Store,
	tokenMaker token.Maker,
	taskDistributor worker.TaskDistributor,
) UseCase {
	return UseCase{
		User:        NewUserUseCase(config, store, tokenMaker, taskDistributor),
		Token:       NewTokenUseCase(config, store, tokenMaker),
		Account:     NewAccountUseCase(store),
		Transaction: NewTransactionUseCase(store),
	}
}

type UserUseCase interface {
	Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, req UpdateUserRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error)
}

type TokenUseCase interface {
	RenewAccessToken(ctx context.Context, req RenewAccessTokenRequest) (*RenewAccessTokenResponse, error)
}

type AccountUseCase interface {
	Create(ctx context.Context, req CreateAccountRequest) (*AccountResponse, error)
	Get(ctx context.Context, req GetAccountRequest) (*AccountResponse, error)
	List(ctx context.Context, req ListAccountRequest) ([]AccountResponse, error)
}

type TransactionUseCase interface {
	Transfer(ctx context.Context, req TransferRequest) (*TransferResult, error)
}
