package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/simplebank/internal/ucase"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	accessToken, err := server.ucase.Token.RenewAccessToken(ctx, ucase.RenewAccessTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken.Token,
		AccessTokenExpiresAt: accessToken.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
