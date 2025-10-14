package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/simplebank/internal/ucase"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	accessToken, err := s.ucase.Token.Renew(ctx, ucase.RenewRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken.Token,
		AccessTokenExpiresAt: accessToken.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
