package http

import (
	"backend-layout/internal/domain"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthHanlder struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(e *echo.Group, au domain.AuthUsecase) {
	handler := &AuthHanlder{
		authUsecase: au,
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
		log.Err(err).Ctx(ctx).
			Str("usecase", "Login").
			Msg("failed to login")

		return err
	}

	return c.JSON(http.StatusOK, &res)
}
