package ucase

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/ymanshur/simplebank/config"
	db "github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/pkg/token"
)

type tokenUcase struct {
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewTokenUseCase(config config.Config, store db.Store, tokenMaker token.Maker) TokenUseCase {
	return &tokenUcase{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RenewAccessTokenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken, validation.Required),
	)
}

type RenewAccessTokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

func (u *tokenUcase) RenewAccessToken(ctx context.Context, req RenewAccessTokenRequest) (*RenewAccessTokenResponse, error) {
	if err := validation.Validate(&req); err != nil {
		return nil, err
	}

	refreshPayload, err := u.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, typex.ErrUnAuthorized(err.Error())
	}

	session, err := u.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, typex.NewErrDataNotFound("token")
		}

		return nil, errors.Wrap(err, "get session")
	}

	if session.IsBlocked {
		return nil, typex.ErrUnAuthorized("session has blocked")
	}

	if session.Username != refreshPayload.Username {
		err = fmt.Errorf("incorrect session user")
		return nil, typex.ErrUnAuthorized("incorrect user session")
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, typex.ErrUnAuthorized("session token mismatched")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, typex.ErrUnAuthorized("session has expired")
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		u.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create token")
	}

	rsp := RenewAccessTokenResponse{
		Token:     accessToken,
		ExpiresAt: accessPayload.ExpiredAt,
	}
	return &rsp, nil
}
