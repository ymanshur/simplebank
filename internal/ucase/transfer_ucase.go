package ucase

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/validator"
)

type transactionUcase struct {
	store db.Store
}

func NewTransactionUseCase(store db.Store) TransactionUseCase {
	return &transactionUcase{
		store: store,
	}
}

type TransferRequest struct {
	Auth          AuthRequest
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (r TransferRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Auth),
		validation.Field(&r.FromAccountID, validation.Required, validator.ValidID()),
		validation.Field(&r.ToAccountID, validation.Required, validator.ValidID()),
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
		validation.Field(&r.Currency, validation.Required, validation.By(validator.ValidCurrency)),
	)
}

type TransferResponse struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	CreatedAt     time.Time
}

type EntryResponse struct {
	ID        int64
	AccountID int64
	Amount    int64
	CreatedAt pgtype.Timestamptz
}

type TransferResult struct {
	Transfer    TransferResponse
	FromAccount AccountResponse
	ToAccount   AccountResponse
	FromEntry   EntryResponse
	ToEntry     EntryResponse
}

func (u *transactionUcase) Transfer(ctx context.Context, req TransferRequest) (*TransferResult, error) {
	if err := validation.Validate(&req); err != nil {
		return nil, err
	}

	fromAccount, err := u.validAccount(ctx, req.FromAccountID, req.Currency)
	if err != nil {
		return nil, err
	}

	if fromAccount.Owner != req.Auth.Username {
		return nil, typex.ErrForbidden("from account doesn't belong to the authenticated user")
	}

	_, err = u.validAccount(ctx, req.ToAccountID, req.Currency)
	if err != nil {
		return nil, err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := u.store.TransferTx(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "transfer")
	}

	rsp := convertTransferResult(result)
	return &rsp, nil
}

func (u *transactionUcase) validAccount(ctx context.Context, accountID int64, currency string) (db.Account, error) {
	account, err := u.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return account, typex.NewErrDataNotFound("account")
		}

		return account, errors.Wrap(err, "get account")
	}

	if account.Currency != currency {
		return account, typex.ErrUnProcessableEnity(fmt.Sprintf("account currency mismatch: %s vs %s", account.Currency, currency))
	}

	return account, nil
}

func convertTransfer(transfer db.Transfer) TransferResponse {
	return TransferResponse{
		ID:            transfer.ID,
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		CreatedAt:     transfer.CreatedAt,
	}
}

func convertEntry(entry db.Entry) EntryResponse {
	return EntryResponse{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		Amount:    entry.Amount,
		CreatedAt: entry.CreatedAt,
	}
}

func convertTransferResult(result db.TransferTxResult) TransferResult {
	return TransferResult{
		Transfer:    convertTransfer(result.Transfer),
		FromAccount: convertAccount(result.FromAccount),
		ToAccount:   convertAccount(result.ToAccount),
		FromEntry:   convertEntry(result.FromEntry),
		ToEntry:     convertEntry(result.ToEntry),
	}
}
