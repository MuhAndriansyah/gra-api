package domain

import (
	"context"
	"time"
)

type Cart struct {
	Id        int64
	UserID    int64
	BookID    int64
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartDetail struct {
	Id          int64
	UserID      int64
	BookID      int64
	BookTitle   string
	ImageUrl    string
	IsAvailable bool
}

type CartResponse struct {
	Id          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	BookID      int64  `json:"book_id"`
	BookTitle   string `json:"book_title"`
	ImageUrl    string `json:"image_url"`
	IsAvailable bool   `json:"is_available"`
}

type StoreCartRequest struct {
	BookID int64 `json:"book_id" validate:"required"`
	UserID int64 `json:"user_id" validate:"required"`
}

func CartDetailToResponse(cartDetail *CartDetail) CartResponse {
	return CartResponse{
		Id:          cartDetail.Id,
		UserID:      cartDetail.UserID,
		BookID:      cartDetail.BookID,
		BookTitle:   cartDetail.BookTitle,
		ImageUrl:    cartDetail.ImageUrl,
		IsAvailable: cartDetail.IsAvailable,
	}
}

type CartRepository interface {
	AddToCart(ctx context.Context, cart *Cart) (id int64, err error)
	DeleteItem(ctx context.Context, cartID, userID int64) error
	CartDetails(ctx context.Context, userID int64) (carts []CartDetail, err error)
	IsItemInCart(ctx context.Context, bookID, userID int64) (bool, error)
	HasStock(ctx context.Context, bookID int64) (bool, error)
}

type CartUsecase interface {
	StoreCart(ctx context.Context, input *StoreCartRequest) (int64, error)
	CartDetails(ctx context.Context, userID int64) (carts []CartResponse, err error)
	DeleteItem(ctx context.Context, cartID, userID int64) error
}
