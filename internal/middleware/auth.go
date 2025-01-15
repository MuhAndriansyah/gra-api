package middleware

import (
	"backend-layout/internal/adapter/jwt"
	"backend-layout/internal/httpcontext"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JWTAuthenticator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			if tokenString == "" || len(tokenString) <= len("Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			tokenString = tokenString[len("Bearer "):]

			claims, err := jwt.ValidateJWT(tokenString)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			c.Set(httpcontext.UserKey, claims)
			return next(c)
		}
	}
}
