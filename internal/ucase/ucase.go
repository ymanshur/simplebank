package ucase

import (
	"github.com/ymanshur/simplebank/config"
	db "github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/worker"
)

type Ucase struct {
	User        UserUcase
	Token       TokenUcase
	Account     AccountUcase
	Transaction TransactionUcase
}

func NewUcase(
	config config.Config,
	repo db.Repo,
	tokenMaker token.Maker,
	taskDistributor worker.TaskDistributor,
) Ucase {
	return Ucase{
		User:        NewUserUcase(config, repo, tokenMaker, taskDistributor),
		Token:       NewTokenUcase(config, repo, tokenMaker),
		Account:     NewAccountUcase(repo),
		Transaction: NewTransactionUcase(repo),
	}
}
