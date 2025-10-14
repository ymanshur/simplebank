package repo

import (
	"context"

	db "github.com/ymanshur/simplebank/db/sqlc"
)

type CreateUserTxParams struct {
	db.CreateUserParams
	AfterCreate func(user db.User) error
}

type CreateUserTxResult struct {
	User db.User
}

func (r *repoQuery) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := r.execTx(ctx, func(q *db.Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
