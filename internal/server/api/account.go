package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
)

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type AccountResponse struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := s.ucase.Account.Create(ctx, ucase.CreateAccountRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
		Currency: req.Currency,
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := convertAccount(account)
	ctx.JSON(http.StatusOK, rsp)
}

type GetAccountRequest struct {
	ID int64 `uri:"id"`
}

func (s *Server) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := s.ucase.Account.Get(ctx, ucase.GetAccountRequest{
		ID: req.ID,
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := convertAccount(account)
	ctx.JSON(http.StatusOK, rsp)
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id"`
	PageSize int32 `form:"page_size"`
}

func (s *Server) ListAccounts(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := s.ucase.Account.List(ctx, ucase.ListAccountRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
		Page: ucase.PagingRequest{
			ID:   common.PageID(req.PageID),
			Size: common.PageSize(req.PageSize),
		},
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := convertAccounts(accounts)
	ctx.JSON(http.StatusOK, rsp)
}
