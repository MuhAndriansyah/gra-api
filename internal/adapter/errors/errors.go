package errors

import (
	"backend-layout/helper"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ErrorHandler ...
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Something went wrong"
	var errorDetails any

	// Prioritaskan baseError
	var baseErr baseError
	if errors.As(err, &baseErr) {
		code = baseErr.code
		message = baseErr.message
	} else if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		message = fmt.Sprintf("%v", httpError.Message)
		errorDetails = httpError.Internal
	} else if vErr, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		message = "Validation error"
		if vld, ok := c.Echo().Validator.(*helper.Validator); ok {
			errorDetails = vld.TranslateError(vErr)
		}
	}

	// Log error yang tidak terduga
	if code >= 500 {
		log.Error().Err(err).Msg("unexpected error")
	}

	c.JSON(code, ErrorResponse{
		Message: message,
		Errors:  errorDetails,
	})
}

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

type baseError struct {
	code    int
	message string
}

func newBaseError(code int, msg string) baseError {
	return baseError{
		code:    code,
		message: msg,
	}
}

func (err baseError) Error() string {
	return strings.ToLower(err.message)
}

func NewNotFoundError(message string) baseError {
	return newBaseError(http.StatusNotFound, message)
}

func NewForbiddenError(message string) baseError {
	return newBaseError(http.StatusForbidden, message)
}

func NewBadRequestError(message string) baseError {
	return newBaseError(http.StatusBadRequest, message)
}

func NewConflictError(message string) baseError {
	return newBaseError(http.StatusConflict, message)
}

func NewUnauthorized(message string) baseError {
	return newBaseError(http.StatusUnauthorized, message)
}

func NewInternalServerError(message string) baseError {
	return newBaseError(http.StatusInternalServerError, message)
}
