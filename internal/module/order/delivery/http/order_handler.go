package http

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"net/http"
	"strconv"

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
	r.GET("/orders", h.Index)
	r.GET("/orders/:id", h.Show)

}

func (h *OrderHandler) Index(c echo.Context) error {
	user, ok := httpcontext.GetUserJWT(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	ctx := c.Request().Context()
	resp, err := h.orderUsecase.GetOrderByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
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

func (h *OrderHandler) Show(c echo.Context) error {
	user, ok := httpcontext.GetUserJWT(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID format")
	}

	ctx := c.Request().Context()
	resp, err := h.orderUsecase.GetOrderDetails(ctx, id, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
