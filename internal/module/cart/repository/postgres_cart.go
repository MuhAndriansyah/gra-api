package repository

import (
	"backend-layout/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresCartRepository struct {
	conn *pgxpool.Pool
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

// HasStock implements domain.CartRepository.
func (p *postgresCartRepository) HasStock(ctx context.Context, bookID int64) (bool, error) {
	query := `SELECT in_stock
			  FROM books
			  WHERE id = $1;`
	var stock int64
	err := p.conn.QueryRow(ctx, query, bookID).Scan(&stock)

	if err != nil {
		return false, err
	}

	return stock > 0, nil
}

// AddToCart implements domain.CartRepository.
func (p *postgresCartRepository) AddToCart(ctx context.Context, cart *domain.Cart) (id int64, err error) {
	query := `INSERT INTO carts (
				user_id, book_id, quantity
			  ) VALUES (
			   $1, $2, $3
			  ) RETURNING id;`

	err = p.conn.QueryRow(ctx, query, cart.UserID, cart.BookID, cart.Quantity).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("could not add item to cart: %w", err)
	}

	return
}

// CartDetails implements domain.CartRepository.
func (p *postgresCartRepository) CartDetails(ctx context.Context, userID int64) (carts []domain.CartDetail, err error) {
	query := `SELECT
				c.id,
  				c.user_id,
  				c.book_id,
  				b.title,
  				b.image_url,
				CASE
					WHEN B.in_stock > 0 THEN TRUE
					ELSE FALSE
				END AS is_available
			  FROM 
			  	carts c
				JOIN books b on c.book_id = b.id
			  WHERE c.user_id = $1;`

	rows, err := p.conn.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	result := make([]domain.CartDetail, 0)

	for rows.Next() {
		c := domain.CartDetail{}

		err := rows.Scan(&c.Id, &c.UserID, &c.BookID, &c.BookTitle, &c.ImageUrl, &c.IsAvailable)

		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, nil
}

// DeleteItem implements domain.CartRepository.
func (p *postgresCartRepository) DeleteItem(ctx context.Context, cartID, userID int64) error {
	query := `DELETE FROM carts WHERE id = $1 and user_id = $2;`

	row, err := p.conn.Exec(ctx, query, cartID, userID)

	if err != nil {
		return err
	}

	rowAffected := row.RowsAffected()

	if rowAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// IsItemInCart implements domain.CartRepository.
func (p *postgresCartRepository) IsItemInCart(ctx context.Context, bookID int64, userID int64) (bool, error) {
	query := `SELECT count(*) FROM carts WHERE book_id = $1 and user_id = $2;`
	var count int
	err := p.conn.QueryRow(ctx, query, bookID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewCartRepository(conn *pgxpool.Pool) domain.CartRepository {
	return &postgresCartRepository{
		conn: conn,
	}
}
