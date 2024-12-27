package api

import (
	"backend-layout/internal/config"
	"backend-layout/internal/middleware"
	"backend-layout/internal/tasks"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type APIServer struct {
	Pool            *pgxpool.Pool
	TaskDistributor tasks.TaskDistributor
	// S3              storage.Uploader
	Conf   *config.Config
	Logger zerolog.Logger
}

func NewAPIServer(pool *pgxpool.Pool, taskDistributor tasks.TaskDistributor, conf *config.Config, logger zerolog.Logger) *APIServer {
	return &APIServer{
		Pool:            pool,
		TaskDistributor: taskDistributor,
		// S3:              s3,
		Conf:   conf,
		Logger: logger,
	}
}

func (s *APIServer) Run(ctx context.Context) error {
	e := echo.New()

	v1 := e.Group("/v1")
	r := v1.Group("")
	p := v1.Group("/public")

	p.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"status": "true",
		})
	})

	r.Use(middleware.JWTAuthenticator())

	go func() {
		<-ctx.Done()
		s.Logger.Info().Msg("shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			s.Logger.Error().Err(err).Msg("server shutdown error")
		}
	}()

	return e.Start(":3000")
}
