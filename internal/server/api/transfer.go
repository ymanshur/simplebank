package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type TransferResponse struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type EntryResponse struct {
	ID        int64              `json:"id"`
	AccountID int64              `json:"account_id"`
	Amount    int64              `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type TransferResult struct {
	Transfer    TransferResponse `json:"transfer"`
	FromAccount AccountResponse  `json:"from_account"`
	ToAccount   AccountResponse  `json:"to_account"`
	FromEntry   EntryResponse    `json:"from_entry"`
	ToEntry     EntryResponse    `json:"to_entry"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	result, err := s.ucase.Transaction.Transfer(ctx, ucase.TransferRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := convertTransferResult(result)
	ctx.JSON(http.StatusOK, rsp)
}
