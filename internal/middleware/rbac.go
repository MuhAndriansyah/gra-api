package middleware

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RBACMiddleware struct {
	rbacService domain.RBACUsecase
}

func NewRBACMiddleware(rbacService domain.RBACUsecase) *RBACMiddleware {
	return &RBACMiddleware{rbacService}
}

func (r *RBACMiddleware) RequiredPermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			au, ok := httpcontext.GetUserJWT(c)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
			}

			hasPermission, err := r.rbacService.CheckUserHasPermission(c.Request().Context(), au.ID, permission)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if !hasPermission {
				return echo.NewHTTPError(http.StatusForbidden, "You don't have permission")
			}

			return next(c)
		}
	}
}

func (r *RBACMiddleware) RequiredRoles(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			au, ok := httpcontext.GetUserJWT(c)

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
			}

			hasRole, err := r.rbacService.CheskUserHasRole(c.Request().Context(), au.ID, role)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "You don't have role")
			}

			return next(c)
		}
	}
}
