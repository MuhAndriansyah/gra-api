package usecase

import (
	"backend-layout/helper"
	baseErr "backend-layout/internal/adapter/errors"
	"backend-layout/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type OrderUsecase struct {
	orderRepo domain.OrderRepository
}

// CreateOrder implements domain.OrderUsecase.
func (o *OrderUsecase) CreateOrder(ctx context.Context, userID int64) (orderResp domain.OrderResponse, err error) {

	tx, err := o.orderRepo.GetTx(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("layer", "usecase").
			Int64("userID", userID).
			Msg("failed to begin transaction")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
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

	items, err := o.orderRepo.GetCartItems(ctx, userID, tx)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to get cart items")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to get items")
	}

	if len(items) == 0 {
		return domain.OrderResponse{}, baseErr.NewBadRequestError("empty cart")
	}

	if len(items) > 3 {
		return domain.OrderResponse{}, baseErr.NewBadRequestError("maximum three items")
	}

	orderNumberStr := generateOrderNumber()

	order := domain.Order{
		UserId:        userID,
		OrderNumber:   orderNumberStr,
		PaymentStatus: "Pending",
	}

	id, err := o.orderRepo.SaveOrder(ctx, &order, tx)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to save item")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
	}

	err = o.orderRepo.SaveOrderDetail(ctx, tx, items, id, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to save order detail")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
	}

	err = o.orderRepo.ClearCart(ctx, userID, tx)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to clear cart")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to commit transaction")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to commit transaction")
	}

	return domain.OrderResponse{
		Id:            id,
		OrderNumber:   orderNumberStr,
		PaymentStatus: "Pending",
		CreatedAt:     time.Now(),
	}, nil

}

func generateOrderNumber() string {
	str, _ := helper.GenerateRandomNumberString(10)

	return "ORD-" + strings.ToUpper(str)
}

func NewOrderUsecase(orderRepo domain.OrderRepository) domain.OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo}
}
