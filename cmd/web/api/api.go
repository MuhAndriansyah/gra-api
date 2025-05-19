package api

import (
	"backend-layout/helper"
	errorHandler "backend-layout/internal/adapter/errors"
	"backend-layout/internal/adapter/oauth"
	paymentgateway "backend-layout/internal/adapter/payment_gateway"
	"backend-layout/internal/config"
	"backend-layout/internal/middleware"
	authHttpDelivery "backend-layout/internal/module/auth/delivery/http"
	_authUsecase "backend-layout/internal/module/auth/usecase"
	bookHttpDelivery "backend-layout/internal/module/book/delivery/http"
	_bookRepository "backend-layout/internal/module/book/repository"
	_bookUsecase "backend-layout/internal/module/book/usecase"
	cartHttpDelivery "backend-layout/internal/module/cart/delivery/http"
	_cartReposiotry "backend-layout/internal/module/cart/repository"
	_cartUsecase "backend-layout/internal/module/cart/usecase"
	_rbacReposiotry "backend-layout/internal/module/rbac/repository"
	_rbacUsecase "backend-layout/internal/module/rbac/usecase"
	userHttpDelivery "backend-layout/internal/module/user/delivery/http"
	_userRepository "backend-layout/internal/module/user/repository"
	_userUsecase "backend-layout/internal/module/user/usecase"
	"fmt"
	"time"

	orderHttpDelivery "backend-layout/internal/module/order/delivery/http"
	_orderRepository "backend-layout/internal/module/order/repository"
	_orderUsecase "backend-layout/internal/module/order/usecase"

	_paymentRepository "backend-layout/internal/module/payment/repository"
	_paymentUsecase "backend-layout/internal/module/payment/usecase"

	paymentHttpDelivery "backend-layout/internal/module/payment/delivery/http"

	"backend-layout/internal/tasks"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type APIServer struct {
	Pool            *pgxpool.Pool
	TaskDistributor tasks.TaskDistributor
	// S3              storage.Uploader
	Conf           *config.Config
	OAuth          *oauth.Oauth
	rdb            *redis.Client
	MidtransClient *paymentgateway.MidtransClient
}

func NewAPIServer(pool *pgxpool.Pool, taskDistributor tasks.TaskDistributor, conf *config.Config, oauth *oauth.Oauth, rdb *redis.Client, midtransClient *paymentgateway.MidtransClient) *APIServer {
	return &APIServer{
		Pool:            pool,
		TaskDistributor: taskDistributor,
		// S3:              s3,
		Conf:           conf,
		OAuth:          oauth,
		rdb:            rdb,
		MidtransClient: midtransClient,
	}
}

var HttpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests",
}, []string{"handler", "method", "code"})

var HttpRequestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "Histogram of response time for handler",
	Buckets: prometheus.DefBuckets,
}, []string{"handler", "method", "code"})

func (s *APIServer) Run(ctx context.Context) error {
	e := echo.New()

	e.Validator = helper.NewValidator()
	e.Use(middleware.CorrelationIDMiddleware)
	e.Use(PrometheusMetricsMiddleware)
	e.HTTPErrorHandler = errorHandler.CustomHTTPErrorHandler

	v1 := e.Group("/api/v1")
	r := v1.Group("")
	p := v1.Group("/public")

	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestsDuration)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	r.Use(middleware.JWTAuthenticator())

	rbacRepository := _rbacReposiotry.NewRBACRepository(s.Pool)
	rbacUsecase := _rbacUsecase.NewRBACUsecase(rbacRepository)
	middlewareRBAC := middleware.NewRBACMiddleware(rbacUsecase)

	userRepository := _userRepository.NewPostgresUserRepository(s.Pool)
	userUsecase := _userUsecase.NewUserUsecase(userRepository, s.TaskDistributor)
	userHttpDelivery.NewUserHandler(p, r, userUsecase)

	authUsecase := _authUsecase.NewAuthUsecase(userRepository, s.rdb)
	authHttpDelivery.NewAuthHandler(p, authUsecase, s.OAuth, s.rdb)

	bookRepository := _bookRepository.NewPostgresBookRepository(s.Pool)
	bookUsecase := _bookUsecase.NewBookUsecase(bookRepository)
	bookHttpDelivery.NewBookHandler(p, r, bookUsecase, middlewareRBAC)

	cartRepository := _cartReposiotry.NewCartRepository(s.Pool)
	cartUsecase := _cartUsecase.NewCartUsecase(cartRepository)
	cartHttpDelivery.NewCartHandler(r, cartUsecase)

	orderRepository := _orderRepository.NewPostgresOrderRepository(s.Pool)
	orderUsecase := _orderUsecase.NewOrderUsecase(orderRepository)
	orderHttpDelivery.NewOrderHandler(r, orderUsecase)

	paymentRepository := _paymentRepository.NewPostgresPaymentRepository(s.Pool)
	paymentUsecase := _paymentUsecase.NewPaymentUsecase(paymentRepository, orderRepository, s.MidtransClient)
	paymentHttpDelivery.NewPaymentHandler(r, paymentUsecase)

	go func() {
		<-ctx.Done()
		log.Info().Msg("shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
	}()

	return e.Start(":3000")
}

func PrometheusMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		duration := time.Since(start).Seconds()

		HttpRequestsTotal.WithLabelValues(c.Request().URL.Path, c.Request().Method, fmt.Sprintf("%v", c.Response().Status)).Inc()
		HttpRequestsDuration.WithLabelValues(c.Request().URL.Path, c.Request().Method, fmt.Sprintf("%v", c.Response().Status)).Observe(duration)
		return err
	}
}
