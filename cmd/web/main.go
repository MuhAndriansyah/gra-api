package main

import (
	"backend-layout/cmd/web/api"
	"backend-layout/internal/adapter/db"
	"backend-layout/internal/adapter/instrumentation"
	"backend-layout/internal/adapter/worker"
	"backend-layout/internal/config"
	"backend-layout/internal/tasks"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := instrumentation.NewLogger(cfg)
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	dbpool, err := initDatabase(cfg, ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
		return
	}
	defer dbpool.Close()

	redisTaskDistributor, err := initRedisTaskDistributor(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize Redis task distributor")
		return
	}
	defer redisTaskDistributor.Close()

	// s3, err := storage.NewS3Client(cfg.AWS)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("error initializing S3 client")
	// 	return
	// }

	srv := api.NewAPIServer(dbpool, redisTaskDistributor, cfg, logger)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, logger)

	if err := srv.Run(ctx); err != nil {
		if err == http.ErrServerClosed {
			logger.Info().Msg("server gracefully stopped")
		} else {
			logger.Fatal().Err(err).Msg("server error")
		}
	}

	if err := waitGroup.Wait(); err != nil {
		logger.Fatal().Err(err).Msg("error from wait group")
	}
}

func initDatabase(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	dbpool, err := db.InitPgx(cfg, ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %w", err)
	}
	return dbpool, nil
}

func initRedisTaskDistributor(cfg *config.Config) (tasks.TaskDistributor, error) {
	redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)}
	return tasks.NewRedisTaskDistributor(redisOpt), nil
}

func runTaskProcessor(ctx context.Context,
	waitGroup *errgroup.Group, logger zerolog.Logger) {

	taskProcessor := worker.NewTaskProcessor()
	waitGroup.Go(func() error {
		if err := taskProcessor.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start task processor")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		logger.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		logger.Info().Msg("task processor is stopped")

		return nil
	})
}
