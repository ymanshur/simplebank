package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/ymanshur/simplebank/internal/ucase"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	user, err := s.ucase.User.Create(ctx, ucase.CreateUserRequest{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		var validationErrors validation.Errors
		if errors.As(err, &validationErrors) {
			ctx.JSON(http.StatusUnprocessableEntity, responseError(err))
			return
		}

		if errors.Is(err, ucase.ErrUserUniqueViolation) {
			ctx.JSON(http.StatusUnprocessableEntity, responseError(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	rsp := convertUser(user)
	ctx.JSON(http.StatusOK, rsp)
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	login, err := s.ucase.User.Login(ctx, ucase.LoginUserRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp:  ctx.ClientIP(),
	})
	if err != nil {
		code, err := translationError(err)
		ctx.JSON(code, responseError(err))
		return
	}

	rsp := LoginUserResponse{
		SessionID:             login.SessionID,
		AccessToken:           login.AccessToken,
		AccessTokenExpiresAt:  login.AccessTokenExpiresAt,
		RefreshToken:          login.RefreshToken,
		RefreshTokenExpiresAt: login.RefreshTokenExpiresAt,
		User:                  convertUser(&login.User),
	}
	ctx.JSON(http.StatusOK, rsp)
}
