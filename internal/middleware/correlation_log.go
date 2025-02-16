package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const CorrelationIDKey string = "correlation_id"

func CorrelationIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		corellationID := uuid.New().String()

		ctx := context.WithValue(c.Request().Context(), CorrelationIDKey, corellationID)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
