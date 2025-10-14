package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/ymanshur/simplebank/db/sqlc"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        db.User
	VerifyEmail db.VerifyEmail
}

func (r *repoQuery) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

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
