package jwt

import (
	"backend-layout/internal/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnknownClaims = errors.New("unknown claims type")
	ErrTokenInvalid  = errors.New("invalid token")
)

type MyClaims struct {
	jwt.RegisteredClaims
	User User `json:"user"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

var privateKey = []byte(config.LoadAppConfig().JWTPrivateKey)

func Sign(ttl time.Duration, user User) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "",
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		User: user,
	})

	return t.SignedString(privateKey)
}

func ValidateJWT(tokenString string) (*User, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigin method: %v", t.Header["alg"])
		}

		return privateKey, nil
	})

	claims, ok := token.Claims.(*MyClaims)

	if !ok {
		return nil, ErrUnknownClaims
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return &claims.User, nil

}
