package main

import (
	"context"
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	statikFS "github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/ymanshur/simplebank/api"
	"github.com/ymanshur/simplebank/db"
	"os"

	"github.com/rs/zerolog/log"
	sqlc "github.com/ymanshur/simplebank/db/sqlc"
	_ "github.com/ymanshur/simplebank/docs/statik"
	"github.com/ymanshur/simplebank/gapi"
	"github.com/ymanshur/simplebank/pb"
	"github.com/ymanshur/simplebank/pkg/util"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}

	err = db.RunMigration(config.MigrationURL, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run database migration")
	}

	store := sqlc.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store sqlc.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	err = server.Start(config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store sqlc.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")

	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//fs := http.FileServer(http.Dir("./docs/swagger"))
	fs, err := statikFS.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik file system")
	}

	//swaggerHandler := http.StripPrefix("/swagger/", fs)
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(fs))
	mux.Handle("/swagger/", swaggerHandler)

	err = server.StartGateway(config.HTTPServerAddress, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func runGinServer(config util.Config, store sqlc.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTP server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP server")
	}
}
