package domain

import (
	"context"
	"time"
)

type User struct {
	Id                       int64
	Name                     string
	Email                    string
	Password                 string
	Photo                    *string
	EmailVerifyCode          string
	EmailVerifyCodeExpiredAt time.Time
	VerifiedAt               *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type VerifyEmailRequest struct {
	Id         int64  `json:"id" validate:"required"`
	VerifyCode string `json:"verify_code" validate:"required"`
}

type StoreUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Store(ctx context.Context, user *User) error
	GetByEmailVerifyCode(ctx context.Context, verifyCode string, id int64) (*User, error)
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, payload *StoreUserRequest) error
	VerifyEmailCode(ctx context.Context, payload *VerifyEmailRequest) error
}
