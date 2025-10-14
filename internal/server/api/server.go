package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/simplebank/config"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/worker"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config          config.Config
	repo            repo.Repo
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor

	httpServer *http.Server
	router     *gin.Engine

	ucase ucase.Ucase
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Config, repo repo.Repo, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		repo:            repo,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
		ucase:           ucase.NewUcase(config, repo, tokenMaker, taskDistributor),
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	if !s.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	v1Routes := router.Group("/v1")
	v1Routes.POST("/create_user", s.CreateUser)
	v1Routes.POST("/login_user", s.LoginUser)
	v1Routes.POST("/renew_access_token", s.RenewAccessToken)

	v1AuthRoutes := v1Routes.Group("/").Use(AuthMiddleware(s.tokenMaker))

	v1AuthRoutes.POST("/accounts", s.CreateAccount)
	v1AuthRoutes.GET("/accounts/:id", s.GetAccount)
	v1AuthRoutes.GET("/accounts", s.ListAccounts)

	v1AuthRoutes.POST("/transfer", s.CreateTransfer)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	s.httpServer = &http.Server{
		Addr:    address,
		Handler: s.router.Handler(),
	}

	log.Info().Msgf("start HTTP server at %s", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	log.Info().Msg("graceful shutdown HTTP server")

	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Info().Msg("HTTP server is stopped")

	return nil
}
