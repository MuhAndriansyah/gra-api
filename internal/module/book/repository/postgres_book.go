package repository

import (
	"backend-layout/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresBookRepository struct {
	conn *pgxpool.Pool
}

/*
sort=highest_price,lowest_price,terbaru
*/

// Fetch implements domain.BookRepository.
func (p *postgresBookRepository) Fetch(ctx context.Context, params domain.BookQueryParams) ([]domain.Book, error) {
	query, args := buildBookQuery(params)

	rows, err := p.conn.Query(ctx, query, args...)

	if err != nil {
		return
	}

	panic("unimplement")
}

// func (p *postgresBookRepository) count(ctx context.Context, query string, args ...interface{}) (total int64, err error) {
// 	err = p.conn.QueryRow(ctx, query, args...).Scan(&total)

// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	return total, nil
// }

func NewPostgresBookRepository(conn *pgxpool.Pool) domain.BookRepository {
	return &postgresBookRepository{
		conn: conn,
	}
}
