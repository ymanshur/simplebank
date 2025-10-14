package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/ymanshur/simplebank/internal/typex"
)

func responseError(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func translationError(err error) (int, error) {
	errCause := errors.Cause(err)
	switch errCause {
	}

	switch errCause := errCause.(type) {
	case validation.Errors, typex.ErrUnProcessableEnity:
		return http.StatusUnprocessableEntity, errCause
	case typex.ErrDataNotFound:
		return http.StatusNotFound, errCause
	case typex.ErrUnAuthorized:
		return http.StatusUnauthorized, errCause
	case typex.ErrForbidden:
		return http.StatusForbidden, errCause
	default:
		return http.StatusInternalServerError, errCause
	}
}
