package usecase

import (
	"backend-layout/internal/adapter/errors"
	"backend-layout/internal/adapter/jwt"
	"backend-layout/internal/domain"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(ur domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		userRepo: ur,
	}
}

// Login implements domain.AuthUsecase.
func (au *authUsecase) Login(ctx context.Context, payload *domain.LoginRequest) (domain.LoginResponse, error) {
	user, err := au.userRepo.GetByEmail(ctx, payload.Email)

	if err != nil {
		return domain.LoginResponse{}, err
	}

	if user == nil {
		return domain.LoginResponse{}, errors.NewUnauthorized("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		return domain.LoginResponse{}, errors.NewUnauthorized("invalid email or password")
	}
	expiration := time.Minute * 60

	token, err := jwt.Sign(expiration, jwt.User{
		ID:    user.Id,
		Email: user.Email,
	})

	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		AccessToken: token,
	}, nil

}
