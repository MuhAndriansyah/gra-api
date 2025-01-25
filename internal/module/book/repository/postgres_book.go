package repository

import (
	"backend-layout/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrISBNDuplicateEntry      = errors.New("duplicate entry: isbn already exists")
	ErrBookForeignKeyViolation = errors.New("foreign key violation: related record not found")
	ErrBookNotFound            = errors.New("book not found")
)

type postgresBookRepository struct {
	conn *pgxpool.Pool
}

// Delete implements domain.BookRepository.
func (p *postgresBookRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM books WHERE id = $1;"

	row, err := p.conn.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	rowAffected := row.RowsAffected()

	if rowAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}

// Update implements domain.BookRepository.
func (p *postgresBookRepository) Update(ctx context.Context, book *domain.Book) error {
	query := `UPDATE books SET
	           title = $1,
						 slug = $2,
						 author_id = $3,
						 publisher_id = $4,
						 publish_year = $5,
						 total_page = $6,
						 description = $7,
						 isbn = $8,
						 price = $9,
						 updated_at = now()
	          WHERE id = $10;`

	row, err := p.conn.Exec(ctx, query,
		book.Title, book.Slug, book.Author.Id, book.Publisher.Id, book.PublishYear, book.TotalPage, book.Description, book.Isbn, book.Price, book.Id)

	if err != nil {
		return err
	}

	rowAffected := row.RowsAffected()

	if rowAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}

// GetByID implements domain.BookRepository.
func (p *postgresBookRepository) GetByID(ctx context.Context, id int64) (*domain.Book, error) {
	query := `SELECT
						        books.id,
        					  books.title,
										books.slug,
										books.author_id,
										authors.name as author_name,
										books.publisher_id,
										publishers.name as publisher_name,
										books.publish_year,
										books.total_page,
										books.description,
										books.sku,
										books.isbn,
										books.price,
										books.created_at,
										books.updated_at
						FROM books WHERE id=$1
						JOIN authors ON authors.id = books.author_id
            JOIN publishers ON publishers.id = books.publisher_id;`

	var b domain.Book

	err := p.conn.QueryRow(ctx, query, id).Scan(&b.Id, &b.Title, &b.Slug, &b.Author.Id, &b.Author.Name, &b.Publisher.Id, &b.Publisher.Name, &b.PublishYear, &b.TotalPage, &b.Description, &b.Sku, &b.Isbn, &b.Price, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	return &b, nil
}

// Store implements domain.BookRepository.
func (p *postgresBookRepository) Store(ctx context.Context, book *domain.Book) (id int64, err error) {
	query := `
	  		INSERT INTO books (
					title, slug, author_id, publisher_id, publish_year, total_page, description, sku, isbn, price
				) VALUES (
				  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
				) RETURNING id;`

	err = p.conn.QueryRow(ctx, query, book.Title, book.Slug, book.Author.Id, book.Publisher.Id, book.PublishYear, book.TotalPage, book.Description, book.Sku, book.Isbn, book.Price).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23503": // foreign key violation
				return 0, ErrBookForeignKeyViolation
			case "23505": // Unique constraint violation
				return 0, ErrISBNDuplicateEntry
			}
		}

		return 0, fmt.Errorf("failed to insert book: %w", err)
	}

	return
}

// Fetch implements domain.BookRepository.
func (p *postgresBookRepository) Fetch(ctx context.Context, params domain.RequestQueryParams) (books []domain.Book, total int64, err error) {

	total, _ = p.count(ctx, params)

	query, args := buildBookQuery(params)

	rows, err := p.conn.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	result := make([]domain.Book, 0)

	for rows.Next() {
		t := domain.Book{}

		err = rows.Scan(
			&t.Id,
			&t.Title,
			&t.Slug,
			&t.Author.Id,
			&t.Author.Name,
			&t.Publisher.Id,
			&t.Publisher.Name,
			&t.PublishYear,
			&t.TotalPage,
			&t.Description,
			&t.Sku,
			&t.Isbn,
			&t.Price,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		result = append(result, t)
	}

	return result, total, nil
}

func (p *postgresBookRepository) count(ctx context.Context, params domain.RequestQueryParams) (total int64, err error) {
	query, args := buildCountBookQuery(params)

	err = p.conn.QueryRow(ctx, query, args...).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func NewPostgresBookRepository(conn *pgxpool.Pool) domain.BookRepository {
	return &postgresBookRepository{
		conn: conn,
	}
}
