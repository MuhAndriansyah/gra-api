package domain

import (
	"context"
	"time"
)

type Book struct {
	Id          int64
	Title       string
	Slug        string
	AuthorId    Author
	PublisherId Publisher
	PublishYear string
	TotalPage   int
	Description string
	Sku         string
	Isbn        string
	Discount    float64
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type BookResponse struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"string"`
	AuthorName    string    `json:"author_name"`
	PublisherName string    `json:"publisher_name"`
	PublishYear   string    `json:"publish_year"`
	TotalPage     int       `json:"total_page"`
	Description   string    `json:"description"`
	Sku           string    `json:"sku"`
	Isbn          string    `json:"isbn"`
	Discount      float64   `json:"discount"`
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func BookToResponse(b *Book) *BookResponse {
	return &BookResponse{
		Id:            b.Id,
		Title:         b.Title,
		Slug:          b.Slug,
		AuthorName:    b.AuthorId.Name,
		PublisherName: b.PublisherId.Name,
		PublishYear:   b.PublishYear,
		TotalPage:     b.TotalPage,
		Description:   b.Description,
		Sku:           b.Sku,
		Isbn:          b.Isbn,
		Discount:      b.Discount,
		Price:         b.Price,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

type BookQueryParams struct {
	MinPrice int64  `query:"minPrice"`
	MaxPrice int64  `query:"maxPrice"`
	Sort     string `query:"sort"`
	Offset   int64  `query:"offset"`
	Search   string `query:"search"`
	Page     int64  `query:"page"`
	PerPage  int64  `query:"perPage"`
}

type BookRepository interface {
	Fetch(ctx context.Context, params BookQueryParams) ([]Book, error)
}
