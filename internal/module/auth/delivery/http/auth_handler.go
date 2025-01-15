package http

import (
	"backend-layout/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHanlder struct {
	authUsecase domain.AuthUsecase
	logger      zerolog.Logger
}

func NewAuthHandler(e *echo.Group, au domain.AuthUsecase, logger zerolog.Logger) {
	handler := &AuthHanlder{
		authUsecase: au,
		logger:      logger,
	}

	e.POST("/users/login", handler.Login)
}

func (h *AuthHanlder) Login(c echo.Context) (err error) {
	req := new(domain.LoginRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	res, err := h.authUsecase.Login(ctx, req)

	if err != nil {
		h.logger.Err(err).Ctx(ctx).
			Str("usecase", "Login").
			Msg("failed to login")

		return err
	}

	return c.JSON(http.StatusOK, &res)
}
