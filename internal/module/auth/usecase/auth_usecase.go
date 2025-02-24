package usecase

import (
	"backend-layout/helper"
	"backend-layout/internal/adapter/errors"
	"backend-layout/internal/adapter/jwt"
	"backend-layout/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo domain.UserRepository
	rdb      *redis.Client
}

func NewAuthUsecase(ur domain.UserRepository, rdb *redis.Client) domain.AuthUsecase {
	return &authUsecase{
		userRepo: ur,
		rdb:      rdb,
	}
}

// LoginOAuth implements domain.AuthUsecase.
func (au *authUsecase) LoginOAuth(ctx context.Context, req *domain.OAuthLoginRequest) (domain.LoginResponse, error) {
	state := req.State
	code := req.Code

	key := fmt.Sprintf("state:%s", state)

	// get state from redis
	stateFromRedis, err := au.rdb.Get(ctx, key).Result()

	if err != nil {
		return domain.LoginResponse{}, errors.NewUnauthorized("invalid state")
	}

	if state != stateFromRedis {
		return domain.LoginResponse{}, errors.NewUnauthorized("invalid state")
	}

	userInfo, err := fetchGoogleUserInfo(code)

	if err != nil || userInfo == nil {
		return domain.LoginResponse{}, err
	}

	user, err := au.userRepo.GetByEmail(ctx, userInfo.Email)

	if err != nil {
		return domain.LoginResponse{}, err
	}

	if user == nil {

		randomPassword, _ := helper.GenerateRandomNumberString(10)

		password, _ := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)

		user = &domain.User{
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			Password:   string(password),
			VerifiedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
		}

		err := au.userRepo.Store(ctx, user)

		if err != nil {
			return domain.LoginResponse{}, err
		}
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

func fetchGoogleUserInfo(accessToken string) (*struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}, error) {
	// fetch user info from google
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create http request :%w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unecpected status code: %d", resp.StatusCode)
	}

	user := struct {
		ID    string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	return &user, nil
}
