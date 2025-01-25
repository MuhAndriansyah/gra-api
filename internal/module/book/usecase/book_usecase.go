package usecase

import (
	"backend-layout/helper"
	"backend-layout/internal/domain"
	"context"
	"strings"
)

type BookUsecase struct {
	bookRepo domain.BookRepository
}

// Store implements domain.BookUsecase.
func (b *BookUsecase) Store(ctx context.Context, payload *domain.StoreBookRequest) (int64, error) {

	sku, _ := helper.GenerateRandomNumberString(10)

	book := domain.Book{
		Title: payload.Title,
		Slug:  toSlug(payload.Title),
		Author: domain.Author{
			Id: payload.AuthorID,
		},
		Publisher: domain.Publisher{
			Id: payload.PublisherID,
		},
		PublishYear: payload.PublishYear,
		TotalPage:   payload.TotalPage,
		Description: payload.Description,
		Sku:         sku,
		Isbn:        payload.Isbn,
		Price:       payload.Price,
	}

	id, err := b.bookRepo.Store(ctx, &book)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Fetch implements domain.BookUseCase.
func (b *BookUsecase) Fetch(ctx context.Context, params domain.RequestQueryParams) ([]domain.Book, int64, error) {
	books, total, err := b.bookRepo.Fetch(ctx, params)

	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func NewBookUsecase(br domain.BookRepository) domain.BookUsecase {
	return &BookUsecase{
		bookRepo: br,
	}
}

func toSlug(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
