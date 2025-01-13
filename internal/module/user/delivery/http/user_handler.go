package http

import (
	"backend-layout/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
	logger      zerolog.Logger
}

func NewUserHanlder(e *echo.Group, r *echo.Group, uu domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: uu,
	}

	e.POST("/users/register", handler.RegisterUser)
	r.POST("/users/email-verification", handler.VerifyEmail)
}

func (h *UserHandler) RegisterUser(c echo.Context) (err error) {
	u := new(domain.StoreUserRequest)

	if err := c.Bind(u); err != nil {
		return err
	}

	if err := c.Validate(u); err != nil {
		return err
	}
	ctx := c.Request().Context()

	err = h.userUsecase.RegisterUser(ctx, u)

	if err != nil {
		h.logger.Err(err).Msg("error in usecase RegisterUser")
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "successfully registered"})
}

func (h *UserHandler) VerifyEmail(c echo.Context) (err error) {
	u := new(domain.VerifyEmailRequest)

	if err := c.Bind(c); err != nil {
		return err
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	ctx := c.Request().Context()

	err = h.userUsecase.VerifyEmailCode(ctx, u)

	if err != nil {
		h.logger.Err(err).Msg("error in usecase VerifyEmail")
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "email verified"})
}
