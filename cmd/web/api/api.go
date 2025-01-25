package api

import (
	"backend-layout/helper"
	errorHandler "backend-layout/internal/adapter/errors"
	"backend-layout/internal/config"
	"backend-layout/internal/middleware"
	authHttpDelivery "backend-layout/internal/module/auth/delivery/http"
	_authUsecase "backend-layout/internal/module/auth/usecase"
	bookHttpDelivery "backend-layout/internal/module/book/delivery/http"
	_bookRepository "backend-layout/internal/module/book/repository"
	_bookUsecase "backend-layout/internal/module/book/usecase"
	userHttpDelivery "backend-layout/internal/module/user/delivery/http"
	_userRepository "backend-layout/internal/module/user/repository"
	_userUsecase "backend-layout/internal/module/user/usecase"
	"backend-layout/internal/tasks"
	"context"

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

	e.Validator = helper.NewValidator()

	e.HTTPErrorHandler = errorHandler.CustomHTTPErrorHandler

	v1 := e.Group("/api/v1")
	r := v1.Group("")
	p := v1.Group("/public")

	r.Use(middleware.JWTAuthenticator())

	userRepository := _userRepository.NewPostgresUserRepository(s.Pool)
	userUsecase := _userUsecase.NewUserUsecase(userRepository, s.TaskDistributor)
	userHttpDelivery.NewUserHanlder(p, r, userUsecase)

	authUsecase := _authUsecase.NewAuthUsecase(userRepository)
	authHttpDelivery.NewAuthHandler(p, authUsecase, s.Logger)

	bookRepository := _bookRepository.NewPostgresBookRepository(s.Pool)
	bookUsecase := _bookUsecase.NewBookUsecase(bookRepository)
	bookHttpDelivery.NewUserHanlder(p, r, bookUsecase, s.Logger)

	go func() {
		<-ctx.Done()
		s.Logger.Info().Msg("shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			s.Logger.Error().Err(err).Msg("server shutdown error")
		}
	}()

	return e.Start(":3000")
}
