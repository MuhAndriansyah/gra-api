package http

import (
	"backend-layout/helper"
	"backend-layout/internal/adapter/oauth"
	"backend-layout/internal/domain"
	"fmt"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
	OAuth       *oauth.Oauth
	rdb         *redis.Client
}

func NewAuthHandler(e *echo.Group, au domain.AuthUsecase, OAuth *oauth.Oauth, rdb *redis.Client) {
	handler := &AuthHandler{
		authUsecase: au,
		OAuth:       OAuth,
		rdb:         rdb,
	}

	e.POST("/users/login", handler.Login)

	e.GET("/auth/google/login", handler.GoogleLogin)
	e.GET("/auth/google/callback", handler.GoogleCallback)
}

func (h *AuthHandler) Login(c echo.Context) (err error) {
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

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	state, _ := helper.GenerateRandomString()

	ctx := c.Request().Context()

	key := fmt.Sprintf("state:%s", state)

	err := h.rdb.Set(ctx, key, state, 2*time.Minute).Err()

	if err != nil {
		return err
	}

	url, err := h.OAuth.AuthUrlGoogleLogin(state)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"url": url})
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	ctx := c.Request().Context()

	oAuthToken, err := h.OAuth.ExhangeCodeForToken(ctx, code)

	if err != nil {
		return err
	}

	res, err := h.authUsecase.LoginOAuth(ctx, &domain.OAuthLoginRequest{
		State: state,
		Code:  oAuthToken.AccessToken,
	})

	if err != nil {
		log.Err(err).Ctx(ctx).
			Str("usecase", "GoogleCallback").
			Msg("failed to google login")

		return err
	}

	return c.JSON(http.StatusOK, &res)
}
