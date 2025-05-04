package repository

import (
	"backend-layout/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresOrderRepository struct {
	conn *pgxpool.Pool
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

// GetByID implements domain.OrderRepository.
func (p *postgresOrderRepository) GetByIDAndUserID(ctx context.Context, id, userId int64) (*domain.Order, error) {
	query := `SELECT id, order_number, user_id, payment_status, payment_date, payment_method, created_at, updated_at 
	          FROM orders 
			  WHERE id = $1 AND user_id = $2;`

	var o domain.Order
	err := p.conn.QueryRow(ctx, query, id, userId).Scan(&o.Id, &o.OrderNumber, &o.UserId, &o.PaymentStatus, &o.PaymentDate, &o.PaymentMethod, &o.CreatedAt, &o.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return &o, nil
}

// GetPendingOrder implements domain.OrderRepository.
func (p *postgresOrderRepository) GetPendingOrder(ctx context.Context, orderNumber string, userId int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE order_number = $1 AND payment_status = 'Pending' AND user_id = $2);`
	var exists bool
	err := p.conn.QueryRow(ctx, query, orderNumber, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// UpdateBorrowDates implements domain.PaymentRepository.
func (p *postgresOrderRepository) UpdateBorrowDates(ctx context.Context, tx pgx.Tx, orderId int64) error {
	query := `UPDATE order_details
	          SET borrowing_date = timezone('UTC', now()),
			  return_date = timezone('UTC', now() + interval '7 days')
			  WHERE order_id = $1;`

	cmdTag, err := tx.Exec(ctx, query, orderId)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated for order_id %v", orderId)
	}

	return nil
}

// UpdateStock implements domain.OrderRepository.
func (p *postgresOrderRepository) UpdateStock(ctx context.Context, tx pgx.Tx, orderID int64) error {
	query := `UPDATE books b 
			  SET in_stock = in_stock - 1
			  FROM order_details od
			  WHERE od.book_id = b.id AND b.in_stock > 0 AND od.order_id = $1;`

	cmdTag, err := tx.Exec(ctx, query, orderID)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("stock not enough")
	}

	return nil
}

// GetOrderByUserID implements domain.OrderRepository.
func (p *postgresOrderRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]domain.OrderWithDetailCount, error) {
	query := `SELECT 
				o.id AS order_id,
				o.order_number,
				o.user_id,
				o.payment_status,
				o.payment_date,
				o.payment_method,
				COUNT(od.id) AS total_order_details,
				o.created_at
			  FROM orders o
			  LEFT JOIN order_details od ON o.id = od.order_id
			  WHERE o.user_id = $1
			  GROUP BY o.id, o.order_number, o.user_id, o.payment_status, o.payment_date, o.payment_method
			  ORDER BY o.created_at DESC;`

	rows, err := p.conn.Query(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("GetOrderByUserID query failed: %w", err)
	}

	defer rows.Close()

	result := make([]domain.OrderWithDetailCount, 0)

	for rows.Next() {
		orderDetail := domain.OrderWithDetailCount{}

		err = rows.Scan(&orderDetail.Id,
			&orderDetail.OrderNumber,
			&orderDetail.UserId,
			&orderDetail.PaymentStatus,
			&orderDetail.PaymentDate,
			&orderDetail.PaymentMethod,
			&orderDetail.TotalOrderDetail,
			&orderDetail.CreatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, orderDetail)
	}

	return result, nil
}

// GetOrderDetail implements domain.OrderRepository.
func (p *postgresOrderRepository) GetOrderDetailWithBook(ctx context.Context, orderID int64) ([]domain.OrderDetailWithBook, error) {
	query := `SELECT 
				od.id, 
				od.order_id, 
				od.book_id,
				od.borrowing_date,
				od.return_date,
				od.created_at,
				od.updated_at, 
				b.title, 
				b.description, 
				b.total_page, 
				b.publish_year,
				a.name as author_name, 
				p.name as publisher_name  
			  FROM order_details od
			  JOIN books b ON od.book_id = b.id
			  JOIN authors a ON b.author_id = a.id
			  JOIN publishers p ON b.publisher_id = p.id 
			  WHERE od.order_id = $1;`

	rows, err := p.conn.Query(ctx, query, orderID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]domain.OrderDetailWithBook, 0)

	for rows.Next() {
		od := domain.OrderDetailWithBook{}

		err = rows.Scan(
			&od.Id,
			&od.OrderId,
			&od.BookId,
			&od.BorrowingDate,
			&od.ReturnDate,
			&od.CreatedAt,
			&od.UpdatedAt,
			&od.BookTitle,
			&od.Description,
			&od.TotalPage,
			&od.PublishYear,
			&od.AuthorName,
			&od.PublisherName)

		if err != nil {
			return nil, err
		}

		result = append(result, od)
	}

	return result, nil
}

// ClearCart implements domain.OrderRepository.
func (p *postgresOrderRepository) ClearCart(ctx context.Context, tx pgx.Tx, userID int64) error {
	query := `DELETE FROM carts 
			  WHERE user_id = $1 AND id IN (
			 	SELECT c.id
				FROM carts c
				JOIN books b ON c.book_id = b.id 
				WHERE b.in_stock > 0
			  );`

	row, err := tx.Exec(ctx, query, userID)

	if err != nil {
		return err
	}

	rowAffected := row.RowsAffected()

	if rowAffected == 0 {
		return fmt.Errorf("no items in cart")
	}

	return nil
}

// GetCartItems implements domain.OrderRepository.
func (p *postgresOrderRepository) GetCartItems(ctx context.Context, tx pgx.Tx, userID int64) ([]*domain.CartItem, error) {
	query := `SELECT c.book_id
			  FROM carts c
			  JOIN books b ON c.book_id = b.id
			  WHERE c.user_id = $1 and b.in_stock > 0;`

	result := make([]*domain.CartItem, 0)

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetCartItems query failed: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		c := &domain.CartItem{}

		err := rows.Scan(&c.BookId)

		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// SaveOrderDetail implements domain.OrderRepository.
func (p *postgresOrderRepository) SaveOrderDetailsFromCart(ctx context.Context, tx pgx.Tx, items []*domain.CartItem, orderID, userID int64) error {

	newItems := make([]*domain.CartItem, len(items))

	if len(items) == 0 {
		return fmt.Errorf("no items to insert in SaveOrderDetailsFromCart")
	}

	for i, v := range items {
		newItem := *v
		newItem.OrderId = orderID
		newItems[i] = &newItem
	}

	_, err := tx.CopyFrom(ctx, pgx.Identifier{"order_details"}, []string{"book_id", "order_id"}, pgx.CopyFromSlice(len(newItems), func(i int) ([]any, error) {
		item := newItems[i]
		return []any{item.BookId, item.OrderId}, nil
	}))

	if err != nil {
		return err
	}

	return nil
}

// GetTx implements domain.OrderRepository.
func (p *postgresOrderRepository) GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.conn.Begin(ctx)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SaveOrder implements domain.OrderRepository.
func (p *postgresOrderRepository) SaveOrder(ctx context.Context, tx pgx.Tx, order *domain.Order) (id int64, err error) {
	query := `INSERT INTO orders(order_number, user_id, payment_status) VALUES ($1, $2, $3) RETURNING id;`

	err = tx.QueryRow(ctx, query, order.OrderNumber, order.UserId, order.PaymentStatus).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	return
}

func NewPostgresOrderRepository(conn *pgxpool.Pool) domain.OrderRepository {
	return &postgresOrderRepository{conn: conn}
}
