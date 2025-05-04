package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PaymentRequest struct {
	OrderNumber string `json:"order_number" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
	UserId      int64  `json:"user_id" validate:"required"`
	Email       string `json:"email" validate:"required"`
}

type PaymentResponse struct {
	Amount     int64  `json:"amount"`
	Token      string `json:"token"`
	PaymentURL string `json:"payment_url"`
}

type PaymentStatusResponse struct {
	PaymentType   string `json:"payment_type"`
	PaymentStatus string `json:"payment_status"`
}

type PaymentStatusRequest struct {
	OrderId int64
	UserId  int64
}

type PaymentRepository interface {
	GetTx(ctx context.Context) (pgx.Tx, error)
	ProcessPayment(ctx context.Context, tx pgx.Tx, userId int64, paymentMethod, orderNumber string) error
	SetPaymentFailed(ctx context.Context, tx pgx.Tx, orderId int64) error
}

type PaymentUsecase interface {
	CreatePayment(ctx context.Context, input *PaymentRequest) (PaymentResponse, error)
	CheckPaymentStatus(ctx context.Context, input *PaymentStatusRequest) (PaymentStatusResponse, error)
}
