package gapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/simplebank/config"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/worker"
	pb "github.com/ymanshur/simplebank/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	// UnimplementedSimpleBankServer enable forward compatibility
	pb.UnimplementedSimpleBankServer

	config          config.Config
	store           repo.Repo
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor

	rpcServer     *grpc.Server
	gatewayServer *http.Server

	ucase ucase.UseCase
}

// NewServer creates a new gRPC server.
func NewServer(config config.Config, repo repo.Repo, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           repo,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
		ucase:           ucase.NewUseCase(config, repo, tokenMaker, taskDistributor),
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		log.Error().
			Any("panic", p).
			Bytes("stack", debug.Stack()).
			Msg("recovered from panic")

		return status.Errorf(codes.Internal, "%s", p)
	}
	grpcRecovery := recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))

	server.rpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		GrpcLogger,
		grpcRecovery,
	))
	pb.RegisterSimpleBankServer(server.rpcServer, server)

	// Allows the gRPC client to explore available RPCs on the server
	// as some kind of self server documentation.
	reflection.Register(server.rpcServer)

	return server, nil
}

func (s *Server) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	return s.rpcServer.Serve(listener)
}

func (s *Server) Shutdown() {
	log.Info().Msg("graceful shutdown gRPC server")

	s.rpcServer.GracefulStop()

	log.Info().Msg("gRPC server is stopped")
}

func (s *Server) StartGateway(address string, allowedOrigins []string, mux *http.ServeMux) error {
	handler := HttpLogger(HttpRecovery(mux))
	handlerWithCORS := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedHeaders: []string{"*"},
	}).Handler(handler)

	s.gatewayServer = &http.Server{
		Addr:    address,
		Handler: handlerWithCORS,
	}

	log.Info().Msgf("start HTTP gateway server at %s", s.gatewayServer.Addr)

	return s.gatewayServer.ListenAndServe()
}

func (s *Server) ShutdownGateway() error {
	log.Info().Msg("graceful shutdown HTTP gateway server")

	err := s.gatewayServer.Shutdown(context.Background())
	if err != nil {
		return err
	}

	log.Info().Msg("HTTP gateway server is stopped")

	return nil
}
