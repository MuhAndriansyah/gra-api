package httpcontext

import (
	"backend-layout/internal/adapter/jwt"

	"github.com/labstack/echo/v4"
)

const UserKey = "user"

func GetUserJWT(c echo.Context) (*jwt.User, bool) {
	user, ok := c.Get(UserKey).(*jwt.User)
	return user, ok
}
