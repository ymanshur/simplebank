package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	statikFS "github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/simplebank/config"
	_ "github.com/ymanshur/simplebank/docs/statik"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/server/api"
	"github.com/ymanshur/simplebank/internal/server/gapi"
	"github.com/ymanshur/simplebank/internal/server/worker"
	"github.com/ymanshur/simplebank/pkg/mail"
	pb "github.com/ymanshur/simplebank/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	}

	if !config.Debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}

	runDBMigration(config.DBMigrationURL, config.DBSource)

	store := repo.NewRepo(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr:     config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	wg, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, wg, config, redisOpt, store)
	runGinServer(ctx, wg, config, store, taskDistributor)
	runGrpcServer(ctx, wg, config, store, taskDistributor)
	runGrpcGatewayServer(ctx, wg, config, store, taskDistributor)

	err = wg.Wait()
	if err != nil {
		log.Fatal().Err(err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config config.Config, redisOpt asynq.RedisClientOpt, repo repo.Repo) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	taskProcessor := worker.NewRedisTaskProcessor(config, redisOpt, repo, mailer)

	log.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()

		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()

		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config config.Config, repo repo.Repo, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, repo, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	waitGroup.Go(func() error {
		err = server.Start(config.GRPCServerAddress)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}

			log.Error().Err(err).Msg("cannot start gRPC server")

			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		server.Shutdown()

		return nil
	})
}

func runGrpcGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Config,
	repo repo.Repo,
	taskDistributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, repo, taskDistributor)
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

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs, err := statikFS.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik file system")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(fs))
	mux.Handle("/swagger/", swaggerHandler)

	waitGroup.Go(func() error {
		err = server.StartGateway(config.GRPCGatewayServerAddress, config.AllowedOrigins, mux)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			log.Error().Err(err).Msg("cannot start HTTP gateway server")

			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		err = server.ShutdownGateway()
		if err != nil {
			log.Error().Err(err).Msg("cannot graceful shutdown HTTP gateway server")

			return err
		}

		return nil
	})
}

func runGinServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Config,
	repo repo.Repo,
	taskDistributor worker.TaskDistributor,
) {
	server, err := api.NewServer(config, repo, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTP server")
	}

	waitGroup.Go(func() error {
		err = server.Start(config.HTTPServerAddress)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			log.Error().Err(err).Msg("cannot start HTTP server")

			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()

		err = server.Shutdown()
		if err != nil {
			log.Error().Err(err).Msg("cannot graceful shutdown HTTP server")

			return err
		}

		return nil
	})
}
