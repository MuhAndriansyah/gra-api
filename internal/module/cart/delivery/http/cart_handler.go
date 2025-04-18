package http

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartUsecase domain.CartUsecase
}

func NewCartHandler(r *echo.Group, cartUsecase domain.CartUsecase) {
	handler := &CartHandler{
		cartUsecase: cartUsecase,
	}

	r.POST("/carts", handler.AddToCart)
	r.DELETE("/carts/:id", handler.DeleteItemCart)
	r.GET("/carts", handler.CartDetail)
}

func (h *CartHandler) AddToCart(c echo.Context) error {

	user, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	req := new(domain.StoreCartRequest)

	req.UserID = user.ID

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	
	id, err := h.cartUsecase.StoreCart(ctx, req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "successfully added item", "data": echo.Map{
		"id": id,
	}})
}

func (h *CartHandler) DeleteItemCart(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cart ID")
	}

	ctx := c.Request().Context()

	user, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	err = h.cartUsecase.DeleteItem(ctx, id, user.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "cart deleted succesfully"})
}

func (h *CartHandler) CartDetail(c echo.Context) error {

	user, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	ctx := c.Request().Context()

	carts, err := h.cartUsecase.CartDetails(ctx, user.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, carts)
}
