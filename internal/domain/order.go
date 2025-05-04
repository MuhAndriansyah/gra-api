package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Order struct {
	Id            int64
	OrderNumber   string
	UserId        int64
	PaymentStatus string
	PaymentDate   *time.Time
	PaymentMethod *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderDetail struct {
	Id            int64
	OrderId       int64
	BookId        int64
	BorrowingDate *time.Time
	ReturnDate    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CartItem struct {
	BookId  int64
	OrderId int64
}

type OrderDetailWithBook struct {
	OrderDetail
	BookTitle     string
	Description   string
	TotalPage     int
	PublishYear   int
	AuthorName    string
	PublisherName string
}

type OrderWithDetailCount struct {
	Order
	TotalOrderDetail int64
}

type OrderResponse struct {
	Id               int64      `json:"id"`
	OrderNumber      string     `json:"order_number"`
	PaymentStatus    string     `json:"payment_status"`
	PaymentDate      *time.Time `json:"payment_date"`
	UserId           int64      `json:"user_id"`
	TotalOrderDetail int64      `json:"total_order_detail,omitempty"`
	CreatedAt        time.Time  `json:"created_at,omitempty"`
}

type OrderDetailResponse struct {
	Id            int64      `json:"id"`
	OrderId       int64      `json:"order_id"`
	BookId        int64      `json:"book_id"`
	BorrowingDate *time.Time `json:"borrowing_date"`
	ReturnDate    *time.Time `json:"return_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	BookTitle     string     `json:"book_title"`
	Description   string     `json:"description"`
	PublisherName string     `json:"publisher_name"`
	PublishYear   int        `json:"publish_year"`
	AuthorName    string     `json:"author_name"`
	TotalPage     int        `json:"total_page"`
}

type OrderRepository interface {
	GetTx(ctx context.Context) (pgx.Tx, error)

	SaveOrder(ctx context.Context, tx pgx.Tx, order *Order) (id int64, err error)
	SaveOrderDetailsFromCart(ctx context.Context, tx pgx.Tx, items []*CartItem, orderID, userID int64) error

	GetByIDAndUserID(ctx context.Context, id, userId int64) (*Order, error)
	GetCartItems(ctx context.Context, tx pgx.Tx, userID int64) ([]*CartItem, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]OrderWithDetailCount, error)
	GetOrderDetailWithBook(ctx context.Context, orderID int64) ([]OrderDetailWithBook, error)
	GetPendingOrder(ctx context.Context, orderNumber string, userId int64) (bool, error)

	UpdateStock(ctx context.Context, tx pgx.Tx, orderID int64) error
	UpdateBorrowDates(ctx context.Context, tx pgx.Tx, orderId int64) error
	ClearCart(ctx context.Context, tx pgx.Tx, userID int64) error
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, userID int64) (orderResp OrderResponse, err error)
	GetUserOrderHistory(ctx context.Context, userID int64) ([]OrderResponse, error)
	GetUserOrderDetails(ctx context.Context, orderID, userID int64) ([]OrderDetailResponse, error)
}
