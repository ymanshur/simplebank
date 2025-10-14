package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/ymanshur/simplebank/db/sqlc"
)

type VerifyUserTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyUserTxResult struct {
	User        db.User
	VerifyEmail db.VerifyEmail
}

func (r *repoQuery) VerifyUserTx(ctx context.Context, arg VerifyUserTxParams) (VerifyUserTxResult, error) {
	var result VerifyUserTxResult

	err := r.execTx(ctx, func(q *db.Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, db.UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, db.UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
