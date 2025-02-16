package domain

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Book struct {
	Id           int64
	Title        string
	Slug         string
	Author       Author
	Publisher    Publisher
	PublishYear  int
	TotalPage    int
	Description  string
	Sku          string
	Isbn         string
	Price        float64
	CategoryName string
	CategoryID   []int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type BookResponse struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"string"`
	AuthorName    string    `json:"author_name"`
	PublisherName string    `json:"publisher_name"`
	PublishYear   int       `json:"publish_year"`
	TotalPage     int       `json:"total_page"`
	Description   string    `json:"description"`
	Sku           string    `json:"sku"`
	Isbn          string    `json:"isbn"`
	Price         float64   `json:"price"`
	CategoryName  []string  `json:"category_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func BookToResponse(b *Book) *BookResponse {
	return &BookResponse{
		Id:            b.Id,
		Title:         b.Title,
		Slug:          b.Slug,
		AuthorName:    b.Author.Name,
		PublisherName: b.Publisher.Name,
		PublishYear:   b.PublishYear,
		TotalPage:     b.TotalPage,
		Description:   b.Description,
		Sku:           b.Sku,
		Isbn:          b.Isbn,
		Price:         b.Price,
		CategoryName:  strings.Split(b.CategoryName, ","),
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

type StoreBookRequest struct {
	AuthorID    int64   `json:"author_id" validate:"required"`
	PublisherID int64   `json:"publisher_id" validate:"required"`
	TotalPage   int     `json:"total_page" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Title       string  `json:"title" validate:"required"`
	PublishYear int     `json:"publish_year" validate:"required,min=1000"`
	Description string  `json:"description"`
	Isbn        string  `json:"isbn" validate:"required,isbn"`
	CategoryID  []int64 `json:"category_id" validate:"required,dive"`
}

type UpdateBookRequest struct {
	ID          int64   `json:"id" validate:"required"`
	AuthorID    int64   `json:"author_id" validate:"required"`
	PublisherID int64   `json:"publisher_id" validate:"required"`
	TotalPage   int     `json:"total_page" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Title       string  `json:"title" validate:"required"`
	PublishYear int     `json:"publish_year" validate:"required,min=1000"`
	Description string  `json:"description"`
	Isbn        string  `json:"isbn" validate:"required,isbn"`
	CategoryID  []int64 `json:"category_id" validate:"required"`
}

type BookRepository interface {
	Fetch(ctx context.Context, params RequestQueryParams) (books []Book, total int64, err error)
	GetTx(ctx context.Context) (pgx.Tx, error)
	Store(ctx context.Context, book *Book) (id int64, err error)
	GetByID(ctx context.Context, id int64) (*Book, error)
	Update(ctx context.Context, tx pgx.Tx, book *Book) error
	Delete(ctx context.Context, id int64) error
}

type BookUsecase interface {
	Fetch(ctx context.Context, params RequestQueryParams) ([]Book, int64, error)
	Store(ctx context.Context, input *StoreBookRequest) (int64, error)
	Update(ctx context.Context, input *UpdateBookRequest) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*BookResponse, error)
}
