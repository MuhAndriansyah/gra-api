package usecase

import (
	baseErr "backend-layout/internal/adapter/errors"
	paymentgateway "backend-layout/internal/adapter/payment_gateway"
	"backend-layout/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog/log"
)

type PaymentUsecase struct {
	paymentRepo    domain.PaymentRepository
	orderRepo      domain.OrderRepository
	midtransClient *paymentgateway.MidtransClient
}

// CheckPaymentStatus implements domain.PaymentUsecase.
func (p *PaymentUsecase) CheckPaymentStatus(ctx context.Context, input *domain.PaymentStatusRequest) (domain.PaymentStatusResponse, error) {

	order, err := p.orderRepo.GetByIDAndUserID(ctx, input.OrderId, input.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PaymentStatusResponse{}, baseErr.NewNotFoundError(err.Error())
		}
		return domain.PaymentStatusResponse{}, baseErr.NewInternalServerError(err.Error())
	}

	transactionStatusResp, errCheckTrx := p.midtransClient.Coreapi.CheckTransaction(order.OrderNumber)
	if errCheckTrx != nil {
		return domain.PaymentStatusResponse{}, err
	}

	tx, err := p.paymentRepo.GetTx(ctx)
	if err != nil {
		return domain.PaymentStatusResponse{}, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Err(rbErr).Str("layer", "usecase").Msg("failed to rollback tx")
			}

			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("original error: %w, rollback error: %v", err, rbErr)
			}
		}
	}()

	if transactionStatusResp != nil {
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				if err := p.paymentRepo.ProcessPayment(ctx, tx, input.UserId, transactionStatusResp.PaymentType, order.OrderNumber); err != nil {
					return domain.PaymentStatusResponse{}, echo.NewHTTPError(http.StatusOK, err.Error())
				}

				if err := p.orderRepo.UpdateBorrowDates(ctx, tx, order.Id); err != nil {
					return domain.PaymentStatusResponse{}, echo.NewHTTPError(http.StatusOK, err.Error())
				}
			} else if transactionStatusResp.FraudStatus == "accept" {
				if err := p.paymentRepo.ProcessPayment(ctx, tx, input.UserId, transactionStatusResp.PaymentType, order.OrderNumber); err != nil {
					return domain.PaymentStatusResponse{}, echo.NewHTTPError(http.StatusOK, err.Error())
				}

				if err := p.orderRepo.UpdateBorrowDates(ctx, tx, order.Id); err != nil {
					return domain.PaymentStatusResponse{}, echo.NewHTTPError(http.StatusOK, err.Error())
				}
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			if err := p.paymentRepo.ProcessPayment(ctx, tx, input.UserId, transactionStatusResp.PaymentType, order.OrderNumber); err != nil {
				return domain.PaymentStatusResponse{}, echo.NewHTTPError(http.StatusOK, err.Error())
			}

			if err := p.orderRepo.UpdateBorrowDates(ctx, tx, order.Id); err != nil {
				return domain.PaymentStatusResponse{}, err
			}
		} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			if err := p.paymentRepo.SetPaymentFailed(ctx, tx, order.Id); err != nil {
				return domain.PaymentStatusResponse{}, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", input.UserId).Int64("orderId", input.OrderId).Msg("failed to commit transaction")

		return domain.PaymentStatusResponse{}, baseErr.NewInternalServerError("failed to commit transaction")
	}

	return domain.PaymentStatusResponse{
		PaymentType:   transactionStatusResp.PaymentType,
		PaymentStatus: transactionStatusResp.TransactionStatus,
	}, err
}

// CreatePayment implements domain.PaymentUsecase.
func (p *PaymentUsecase) CreatePayment(ctx context.Context, input *domain.PaymentRequest) (domain.PaymentResponse, error) {

	isOrder, err := p.orderRepo.GetPendingOrder(ctx, input.OrderNumber, input.UserId)

	if err != nil {
		return domain.PaymentResponse{}, err
	}

	if !isOrder {
		return domain.PaymentResponse{}, baseErr.NewNotFoundError("order not found")
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  input.OrderNumber,
			GrossAmt: input.Amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: input.Email,
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}

	snapResp, _ := p.midtransClient.SnapClient.CreateTransaction(snapReq)

	return domain.PaymentResponse{
		Amount:     input.Amount,
		Token:      snapResp.Token,
		PaymentURL: snapResp.RedirectURL,
	}, nil

}

func NewPaymentUsecase(paymentRepo domain.PaymentRepository, orderRepo domain.OrderRepository, midtransClient *paymentgateway.MidtransClient) domain.PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepo,
		midtransClient: midtransClient,
	}
}
