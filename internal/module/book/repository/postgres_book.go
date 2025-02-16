package repository

import (
	"backend-layout/internal/domain"
	"backend-layout/internal/middleware"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	ErrISBNDuplicateEntry      = errors.New("duplicate entry: isbn already exists")
	ErrBookForeignKeyViolation = errors.New("foreign key violation: related record not found")
	ErrBookNotFound            = errors.New("book not found")
)

type postgresBookRepository struct {
	conn *pgxpool.Pool
}

// GetTx implements domain.BookRepository.
func (p *postgresBookRepository) GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.conn.Begin(ctx)

	if err != nil {
		return nil, err
	}

	return tx, nil
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
func (p *postgresBookRepository) Update(ctx context.Context, tx pgx.Tx, book *domain.Book) error {
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

	row, err := tx.Exec(ctx, query,
		book.Title, book.Slug, book.Author.Id, book.Publisher.Id, book.PublishYear, book.TotalPage, book.Description, book.Isbn, book.Price, book.Id)

	if err != nil {
		return err
	}

	rowAffected := row.RowsAffected()

	if rowAffected == 0 {
		return ErrBookNotFound
	}

	if _, err := tx.Exec(ctx, `DELETE FROM book_category WHERE book_id=$1`, book.Id); err != nil {
		return fmt.Errorf("failed to delete existing categories :%w", err)
	}
	
	_, err = tx.CopyFrom(ctx, pgx.Identifier{"book_category"}, []string{"category_id", "book_id"}, pgx.CopyFromSlice(len(book.CategoryID), func(i int) ([]any, error) {
		return []any{book.CategoryID[i], book.Id}, nil
	}))

	if err != nil {
		return err
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
				STRING_AGG(c.name, ',') as category_name,
				books.created_at,
				books.updated_at
			FROM books
			JOIN authors ON authors.id = books.author_id
            JOIN publishers ON publishers.id = books.publisher_id
			JOIN book_category bc ON books.id = bc.book_id
			JOIN categories c ON bc.category_id = c.id
			WHERE books.id=$1
			GROUP BY books.id, books.title, books.slug, books.author_id, 
			authors.name, books.publisher_id, publishers.name, 
			books.publish_year, books.total_page, books.description, 
			books.sku, books.isbn, books.price, books.created_at, books.updated_at;`

	var b domain.Book

	err := p.conn.QueryRow(ctx, query, id).Scan(&b.Id, &b.Title, &b.Slug, &b.Author.Id, &b.Author.Name, &b.Publisher.Id, &b.Publisher.Name, &b.PublishYear, &b.TotalPage, &b.Description, &b.Sku, &b.Isbn, &b.Price, &b.CategoryName, &b.CreatedAt, &b.UpdatedAt)

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

		log.Err(err).
			Str("corellation_id", ctx.Value(middleware.CorrelationIDKey).(string)).
			Str("service", "book-service").
			Str("layer", "repository").
			Msg("Failed to execute query")

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

	_, err = p.conn.CopyFrom(ctx, pgx.Identifier{"book_category"}, []string{"category_id", "book_id"}, pgx.CopyFromSlice(len(book.CategoryID), func(i int) ([]any, error) {
		return []any{book.CategoryID[i], id}, nil
	}))

	if err != nil {
		log.Err(err).
			Str("corellation_id", ctx.Value(middleware.CorrelationIDKey).(string)).
			Str("service", "book-service").
			Str("layer", "repository").
			Msg("Failed to execute query")

		return 0, fmt.Errorf("failed to insert book_category: %w", err)
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
			&t.CategoryName,
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
