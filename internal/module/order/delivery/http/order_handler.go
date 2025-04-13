package http

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(r *echo.Group, ou domain.OrderUsecase) {
	h := &OrderHandler{
		orderUsecase: ou,
	}

	r.POST("/orders", h.Store)

}

func (h *OrderHandler) Store(c echo.Context) error {
	user, ok := httpcontext.GetUserJWT(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	ctx := c.Request().Context()

	resp, err := h.orderUsecase.CreateOrder(ctx, user.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}
