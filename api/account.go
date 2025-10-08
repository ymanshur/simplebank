package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

type accountResponse struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.ucase.Account.Create(ctx, ucase.CreateAccountRequest{
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

type getAccountRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.ucase.Account.Get(ctx, ucase.GetAccountRequest{
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

type listAccountRequest struct {
	PageID   int32 `form:"page_id"`
	PageSize int32 `form:"page_size"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := server.ucase.Account.List(ctx, ucase.ListAccountRequest{
		Auth: ucase.AuthRequest{
			Username: authPayload.Username,
		},
		Page: ucase.PagingRequest{
			ID:   PageID(req.PageID),
			Size: PageSize(req.PageSize),
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

func convertAccount(account *ucase.AccountResponse) accountResponse {
	return accountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func convertAccounts(accounts []ucase.AccountResponse) []accountResponse {
	var res []accountResponse
	for _, account := range accounts {
		res = append(res, convertAccount(&account))
	}
	return res
}
