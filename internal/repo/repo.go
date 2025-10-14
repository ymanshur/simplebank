package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ymanshur/simplebank/db/sqlc"
)

// Repo defines all functions to execute db queries and transactions
type Repo interface {
	db.Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyUserTx(ctx context.Context, arg VerifyUserTxParams) (VerifyUserTxResult, error)
}

// repoQuery provides all functions to execute SQL queries and transactions
type repoQuery struct {
	pool *pgxpool.Pool

	// composition
	*db.Queries
}

// NewRepo creates a new Repo
func NewRepo(pool *pgxpool.Pool) Repo {
	return &repoQuery{
		pool:    pool,
		Queries: db.New(pool),
	}
}

// execTx executes a function within a database transaction
func (r *repoQuery) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
