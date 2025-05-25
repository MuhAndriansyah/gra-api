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
	var baseErr BaseError
	if errors.As(err, &baseErr) {
		code = baseErr.Code
		message = baseErr.Message
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

type BaseError struct {
	Code    int
	Message string
}

func newBaseError(code int, msg string) BaseError {
	return BaseError{
		Code:    code,
		Message: msg,
	}
}

func (err BaseError) Error() string {
	return strings.ToLower(err.Message)
}

func NewNotFoundError(message string) BaseError {
	return newBaseError(http.StatusNotFound, message)
}

func NewForbiddenError(message string) BaseError {
	return newBaseError(http.StatusForbidden, message)
}

func NewBadRequestError(message string) BaseError {
	return newBaseError(http.StatusBadRequest, message)
}

func NewConflictError(message string) BaseError {
	return newBaseError(http.StatusConflict, message)
}

func NewUnauthorized(message string) BaseError {
	return newBaseError(http.StatusUnauthorized, message)
}

func NewInternalServerError(message string) BaseError {
	return newBaseError(http.StatusInternalServerError, message)
}
