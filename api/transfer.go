package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type transferResponse struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type entryResponse struct {
	ID        int64              `json:"id"`
	AccountID int64              `json:"account_id"`
	Amount    int64              `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type transferResult struct {
	Transfer    transferResponse `json:"transfer"`
	FromAccount accountResponse  `json:"from_account"`
	ToAccount   accountResponse  `json:"to_account"`
	FromEntry   entryResponse    `json:"from_entry"`
	ToEntry     entryResponse    `json:"to_entry"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	result, err := server.ucase.Transaction.Transfer(ctx, ucase.TransferRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	rsp := convertTransferResult(result)
	ctx.JSON(http.StatusOK, rsp)
}

func convertTransferResult(result *ucase.TransferResult) *transferResult {
	return &transferResult{
		Transfer: transferResponse{
			ID:            result.Transfer.ID,
			FromAccountID: result.Transfer.FromAccountID,
			ToAccountID:   result.Transfer.ToAccountID,
			Amount:        result.Transfer.Amount,
			CreatedAt:     result.Transfer.CreatedAt,
		},
		FromAccount: accountResponse{
			ID:        result.FromAccount.ID,
			Owner:     result.FromAccount.Owner,
			Balance:   result.FromAccount.Balance,
			Currency:  result.FromAccount.Currency,
			CreatedAt: result.FromAccount.CreatedAt,
		},
		ToAccount: accountResponse{
			ID:        result.ToAccount.ID,
			Owner:     result.ToAccount.Owner,
			Balance:   result.ToAccount.Balance,
			Currency:  result.ToAccount.Currency,
			CreatedAt: result.ToAccount.CreatedAt,
		},
		FromEntry: entryResponse{
			ID:        result.FromEntry.ID,
			AccountID: result.FromEntry.AccountID,
			Amount:    result.FromEntry.Amount,
			CreatedAt: result.FromEntry.CreatedAt,
		},
		ToEntry: entryResponse{
			ID:        result.ToEntry.ID,
			AccountID: result.ToEntry.AccountID,
			Amount:    result.ToEntry.Amount,
			CreatedAt: result.ToEntry.CreatedAt,
		},
	}
}
