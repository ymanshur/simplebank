package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/typex"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
	http       *http.Server
	ucase      ucase.UseCase
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		ucase:      ucase.NewUseCase(config, store, tokenMaker, taskDistributor),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	if !server.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	v1Routes := router.Group("/v1")
	v1Routes.POST("/create_user", server.createUser)
	v1Routes.POST("/login_user", server.loginUser)
	v1Routes.POST("/renew_access_token", server.renewAccessToken)

	v1AuthRoutes := v1Routes.Group("/").Use(authMiddleware(server.tokenMaker))

	v1AuthRoutes.POST("/accounts", server.createAccount)
	v1AuthRoutes.GET("/accounts/:id", server.getAccount)
	v1AuthRoutes.GET("/accounts", server.listAccounts)

	v1AuthRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	server.http = &http.Server{
		Addr:    address,
		Handler: server.router.Handler(),
	}

	log.Info().Msgf("start HTTP server at %s", server.http.Addr)

	return server.http.ListenAndServe()
}

func (server *Server) Shutdown() error {
	log.Info().Msg("graceful shutdown HTTP server")

	err := server.http.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Info().Msg("HTTP server is stopped")

	return nil
}

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
