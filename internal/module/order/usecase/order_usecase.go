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

// GetUserOrderHistory implements domain.OrderUsecase.
func (o *OrderUsecase) GetUserOrderHistory(ctx context.Context, userID int64) ([]domain.OrderResponse, error) {
	orders, err := o.orderRepo.GetOrderByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to get order details by user_id")
		return nil, baseErr.NewInternalServerError("failed to get orders")
	}

	orderResponses := make([]domain.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, domain.OrderResponse{
			Id:               order.Id,
			OrderNumber:      order.OrderNumber,
			PaymentStatus:    order.PaymentStatus,
			PaymentDate:      order.PaymentDate,
			UserId:           order.UserId,
			TotalOrderDetail: order.TotalOrderDetail,
			CreatedAt:        order.CreatedAt,
		})
	}

	return orderResponses, nil
}

// GetUserOrderDetails implements domain.OrderUsecase.
func (o *OrderUsecase) GetUserOrderDetails(ctx context.Context, orderID, userID int64) ([]domain.OrderDetailResponse, error) {

	isOrderExists, err := o.orderRepo.IsOrderOwnedByUser(ctx, orderID, userID)

	if err != nil {
		return nil, fmt.Errorf("check order ownership failed: %w", err)
	}

	if !isOrderExists {
		return nil, baseErr.NewNotFoundError("order not found")
	}

	orderDetailsWithBookInfo, err := o.orderRepo.GetOrderDetailWithBook(ctx, orderID)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Int64("orderID", orderID).Msg("failed to get order details")
		return nil, baseErr.NewInternalServerError("failed to get order details")
	}

	orderDetailResponses := make([]domain.OrderDetailResponse, 0, len(orderDetailsWithBookInfo))
	for _, od := range orderDetailsWithBookInfo {
		orderDetailResponses = append(orderDetailResponses, domain.OrderDetailResponse{
			Id:            od.Id,
			OrderId:       od.OrderId,
			BookId:        od.BookId,
			BorrowingDate: od.BorrowingDate,
			ReturnDate:    od.ReturnDate,
			CreatedAt:     od.CreatedAt,
			UpdatedAt:     od.UpdatedAt,
			BookTitle:     od.BookTitle,
			Description:   od.Description,
			PublisherName: od.PublisherName,
			PublishYear:   od.PublishYear,
			AuthorName:    od.AuthorName,
			TotalPage:     od.TotalPage,
		})

	}

	return orderDetailResponses, nil

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

	items, err := o.orderRepo.GetCartItems(ctx, tx, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to get cart items")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to get items")
	}

	if len(items) == 0 {
		return domain.OrderResponse{}, baseErr.NewBadRequestError("empty cart")
	}

	if len(items) > 3 {
		return domain.OrderResponse{}, baseErr.NewBadRequestError("maximum 3 items")
	}

	orderNumberStr := generateOrderNumber()

	order := domain.Order{
		UserId:        userID,
		OrderNumber:   orderNumberStr,
		PaymentStatus: "Pending",
	}

	id, err := o.orderRepo.SaveOrder(ctx, tx, &order)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to save item")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
	}

	err = o.orderRepo.SaveOrderDetailsFromCart(ctx, tx, items, id, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("userID", userID).Msg("failed to save order detail")

		return domain.OrderResponse{}, baseErr.NewInternalServerError("failed to create order")
	}

	err = o.orderRepo.ClearCart(ctx, tx, userID)
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
