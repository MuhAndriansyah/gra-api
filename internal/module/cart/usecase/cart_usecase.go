package usecase

import (
	baseErr "backend-layout/internal/adapter/errors"
	"backend-layout/internal/domain"
	"backend-layout/internal/module/cart/repository"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type CartUsecase struct {
	cartRepo domain.CartRepository
}

// CartDetails implements domain.CartUsecase.
func (c *CartUsecase) CartDetails(ctx context.Context, userID int64) (carts []domain.CartResponse, err error) {
	cartDetails, err := c.cartRepo.CartDetails(ctx, userID)

	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Int64("user_id", userID).Msg("failed to get cart details")
		return nil, baseErr.NewInternalServerError("failed to get cart details")
	}

	cartResponses := make([]domain.CartResponse, 0, len(cartDetails))

	for _, v := range cartDetails {
		cartResponses = append(cartResponses, domain.CartDetailToResponse(&v))
	}

	return cartResponses, nil
}

// DeleteItem implements domain.CartUsecase.
func (c *CartUsecase) DeleteItem(ctx context.Context, cartID, userID int64) error {

	// Hapus item dari cart
	err := c.cartRepo.DeleteItem(ctx, cartID, userID)

	if err != nil {
		log.Error().
			Err(err).
			Str("layer", "usecase").
			Int64("user_id", userID).
			Int64("cart_id", cartID).
			Msg("failed to delete item from cart")

		if errors.Is(err, repository.ErrRecordNotFound) {
			return baseErr.NewNotFoundError(err.Error())
		} else {
			return baseErr.NewInternalServerError(fmt.Sprintf("failed to delete item from cart: %v", err))

		}

	}

	// Log keberhasilan
	log.Info().
		Str("layer", "usecase").
		Int64("user_id", userID).
		Int64("cart_id", cartID).
		Msg("item successfully deleted from cart")

	return nil
}

// StoreCart implements domain.CartUsecase.
func (c *CartUsecase) StoreCart(ctx context.Context, input *domain.StoreCartRequest) (int64, error) {

	isExists, err := c.cartRepo.IsItemInCart(ctx, input.BookID, input.UserID)

	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").
			Int64("user_id", input.UserID).
			Int64("book_id", input.BookID).
			Msg("failed tp check if item exists in cart")

		return 0, baseErr.NewInternalServerError("failed to add item to cart")
	}

	if isExists {
		return 0, baseErr.NewBadRequestError("item already in cart")
	}

	hasStock, err := c.cartRepo.HasStock(ctx, input.BookID)

	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").
			Int64("user_id", input.UserID).
			Int64("book_id", input.BookID).
			Msg("failed to check stock")

		return 0, baseErr.NewInternalServerError("failed to check stock")

	}

	if !hasStock {
		return 0, baseErr.NewConflictError("item is out of stock")
	}

	id, err := c.cartRepo.AddToCart(ctx, &domain.Cart{
		BookID:   input.BookID,
		UserID:   input.UserID,
		Quantity: 1,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("layer", "usecase").
			Int64("user_id", input.UserID).
			Msg("failed to add item to cart")

		return 0, err
	}

	return id, nil
}

func NewCartUsecase(cartRepo domain.CartRepository) domain.CartUsecase {
	return &CartUsecase{cartRepo: cartRepo}
}
