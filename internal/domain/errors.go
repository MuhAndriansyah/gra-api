package domain

import "errors"

var (
	ErrEmailDuplicate = errors.New("email already used")
)

type ErrResponse struct {
	Message string `json:"message"`
}

func NewErrResponse(err error) *ErrResponse {
	return &ErrResponse{
		Message: err.Error(),
	}
}
