package http

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
	logger      zerolog.Logger
}

func NewUserHandler(e *echo.Group, r *echo.Group, uu domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: uu,
	}

	e.POST("/users/register", handler.RegisterUser)
	r.POST("/users/email-verification", handler.VerifyEmail)
}

func (h *UserHandler) RegisterUser(c echo.Context) (err error) {
	req := new(domain.StoreUserRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	err = h.userUsecase.RegisterUser(ctx, req)

	if err != nil {
		h.logger.Err(err).
			Ctx(ctx).
			Str("usecase", "RegisterUser").
			Msg("failed to create user")

		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "successfully registered"})
}

func (h *UserHandler) VerifyEmail(c echo.Context) (err error) {
	req := new(domain.VerifyEmailRequest)
	au, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	req.Id = au.ID

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err = h.userUsecase.VerifyEmailCode(ctx, req)

	if err != nil {
		h.logger.Err(err).Ctx(ctx).
			Str("usecase", "VerifyEmailCode").
			Msg("failed to verify email")

		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "email verified"})
}
