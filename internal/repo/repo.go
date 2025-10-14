package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo defines all functions to execute db queries and transactions
type Repo interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// repoQuery provides all functions to execute SQL queries and transactions
type repoQuery struct {
	db *pgxpool.Pool

	// composition
	*Queries
}

// NewRepo creates a new Repo
func NewRepo(db *pgxpool.Pool) Repo {
	return &repoQuery{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (r *repoQuery) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
