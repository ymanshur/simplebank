package gapi

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/pb"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	// UnimplementedSimpleBankServer enable forward compatibility
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	rpc             *grpc.Server
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	grpcLogger := grpc.UnaryInterceptor(GrpcLogger)
	server.rpc = grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(server.rpc, server)

	// Allows the gRPC client to explore available RPCs on the server
	// as some kind of self server documentation.
	reflection.Register(server.rpc)

	return server, nil
}

func (server *Server) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	return server.rpc.Serve(listener)
}

func (server *Server) StartGateway(address string, mux *http.ServeMux) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := HttpLogger(mux)
	return http.Serve(listener, handler)
}
