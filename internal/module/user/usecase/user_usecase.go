package usecase

import (
	"backend-layout/helper"
	"backend-layout/internal/adapter/errors"
	"backend-layout/internal/domain"
	"backend-layout/internal/tasks"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo        domain.UserRepository
	taskDistributor tasks.TaskDistributor
}

// RegisterUser implements domain.UserUsecase.
func (u *userUsecase) RegisterUser(ctx context.Context, payload *domain.StoreUserRequest) error {
	user, err := u.userRepo.GetByEmail(ctx, payload.Email)

	if err != nil {
		return err
	}

	if user != nil {
		return errors.NewConflictError(domain.ErrEmailDuplicate.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	verifyCode, err := helper.GenerateRandomNumberString(6)

	if err != nil {
		return err
	}

	err = u.userRepo.Store(ctx, &domain.User{
		Name:                     payload.Name,
		Email:                    payload.Email,
		Password:                 string(hashedPassword),
		EmailVerifyCode:          verifyCode,
		EmailVerifyCodeExpiredAt: time.Now().Add(24 * time.Hour),
	})

	if err != nil {
		return err
	}

	err = u.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &tasks.PayloadSendVerifyEmail{
		Email:      payload.Email,
		Username:   payload.Name,
		VerifyCode: verifyCode,
	})

	if err != nil {
		return err
	}

	return nil
}

// VerifyEmailCode implements domain.UserUsecase.
func (u *userUsecase) VerifyEmailCode(ctx context.Context, payload *domain.VerifyEmailRequest) error {
	user, err := u.userRepo.GetByEmailVerifyCode(ctx, payload.VerifyCode, payload.Id)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.NewNotFoundError("invalid verification code")
	}

	if time.Now().After(user.EmailVerifyCodeExpiredAt) {
		return errors.NewNotFoundError("verification code expired")
	}

	err = u.userRepo.ValidatingEmail(ctx, payload.VerifyCode, payload.Id)

	if err != nil {
		return err
	}

	return nil
}

func NewUserUsecase(ur domain.UserRepository, td tasks.TaskDistributor) domain.UserUsecase {
	return &userUsecase{
		userRepo:        ur,
		taskDistributor: td,
	}
}
