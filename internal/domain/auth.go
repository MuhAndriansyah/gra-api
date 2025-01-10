package domain

import "context"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthUsecase interface {
	Login(ctx context.Context, req *LoginRequest) (LoginResponse, error)
}
