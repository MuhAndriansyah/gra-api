package errors

import (
	"backend-layout/helper"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ErrorHandler ...
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Something went wrong"
	var validateError any

	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
	}

	switch v := err.(type) {
	case *echo.HTTPError:
		code = v.Code
		validateError = v.Message
	case validator.ValidationErrors:
		code = http.StatusBadRequest
		message = "Validation error"

		// Access the validator instance from Echo to get the translator
		if vld, ok := c.Echo().Validator.(*helper.Validator); ok {
			validateError = vld.TranslateError(err)
		} else {
			validateError = "Validation errors"
		}
	case baseError:
		code = v.code
		message = v.Error()
	default:
		if err != nil {
			message = err.Error()
		}
	}

	errResponse := ErrorResponse{
		Message: message,
		Errors:  validateError,
	}

	c.JSON(code, errResponse)
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
