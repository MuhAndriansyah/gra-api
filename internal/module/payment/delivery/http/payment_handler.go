package http

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/httpcontext"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	paymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(r *echo.Group, pu domain.PaymentUsecase) {
	handler := &PaymentHandler{paymentUsecase: pu}

	r.POST("/payment", handler.CreatePayment)
	r.GET("/payment/status/:order_id", handler.PaymentStatus)
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	req := new(domain.PaymentRequest)
	req.Amount = 100000

	userJWT, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	if userJWT == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	req.UserId = userJWT.ID
	req.Email = userJWT.Email

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	payment, err := h.paymentUsecase.CreatePayment(ctx, req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, payment)

}

func (h *PaymentHandler) PaymentStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}

	user, ok := httpcontext.GetUserJWT(c)

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized access. Please log in to continue.")
	}

	req := new(domain.PaymentStatusRequest)
	req.OrderId = id
	req.UserId = user.ID
	fmt.Println(req)
	ctx := c.Request().Context()

	resp, err := h.paymentUsecase.CheckPaymentStatus(ctx, req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
