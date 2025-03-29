package usecase

import (
	"backend-layout/helper"
	baseErr "backend-layout/internal/adapter/errors"
	"backend-layout/internal/domain"
	"backend-layout/internal/middleware"
	"backend-layout/internal/module/book/repository"
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type BookUsecase struct {
	bookRepo domain.BookRepository
}

// Delete implements domain.BookUsecase.
func (b *BookUsecase) Delete(ctx context.Context, id int64) error {
	err := b.bookRepo.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return baseErr.NewNotFoundError("book not found")
		}

		log.Error().Err(err).Str("layer", "usecase").Msg("failed to delete book")

		return baseErr.NewInternalServerError("failed to delete book")
	}

	return nil
}

// Get implements domain.BookUsecase.
func (b *BookUsecase) Get(ctx context.Context, id int64) (domain.BookResponse, error) {
	book, err := b.bookRepo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			log.Warn().Int64("book_id", id).Str("layer", "usecase").Msg("book not found")
			return domain.BookResponse{}, baseErr.NewNotFoundError("book not found")
		}

		log.Error().Err(err).Str("layer", "usecase").Int64("book_id", id).Msg("failed to get book")

		return domain.BookResponse{}, baseErr.NewInternalServerError("failed to get book")
	}

	return domain.BookToResponse(book), nil
}

// Update implements domain.BookUsecase.
func (b *BookUsecase) Update(ctx context.Context, input *domain.UpdateBookRequest) error {

	book := domain.Book{
		Id:    input.ID,
		Title: input.Title,
		Slug:  toSlug(input.Title),
		Author: domain.Author{
			Id: input.AuthorID,
		},
		Publisher: domain.Publisher{
			Id: input.PublisherID,
		},
		PublishYear: input.PublishYear,
		TotalPage:   input.TotalPage,
		Description: input.Description,
		Isbn:        input.Isbn,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
	}

	tx, err := b.bookRepo.GetTx(ctx)

	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Msg("failed to begin transaction")

		return baseErr.NewInternalServerError("failed to update book")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if err := b.bookRepo.Update(ctx, tx, &book); err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return baseErr.NewNotFoundError("book not found")
		}

		log.Error().Err(err).Str("layer", "usecase").Msg("failed to update book")

		return baseErr.NewInternalServerError("failed to update book")
	}

	return nil
}

// Store implements domain.BookUsecase.
func (b *BookUsecase) Store(ctx context.Context, input *domain.StoreBookRequest) (int64, error) {

	sku, _ := helper.GenerateRandomNumberString(10)

	book := domain.Book{
		Title: input.Title,
		Slug:  toSlug(input.Title),
		Author: domain.Author{
			Id: input.AuthorID,
		},
		Publisher: domain.Publisher{
			Id: input.PublisherID,
		},
		PublishYear: input.PublishYear,
		TotalPage:   input.TotalPage,
		Description: input.Description,
		Sku:         sku,
		Isbn:        input.Isbn,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
	}

	id, err := b.bookRepo.Store(ctx, &book)

	if err != nil {

		if errors.Is(err, repository.ErrISBNDuplicateEntry) {
			return 0, baseErr.NewNotFoundError("isbn already exist")
		}

		log.Err(err).
			Str("corellation_id", ctx.Value(middleware.CorrelationIDKey).(string)).
			Str("service", "book-service").
			Str("layer", "usecase").
			Str("func", "store").
			Any("parameters", book).Msg("Failed to create book in repository")

		return 0, baseErr.NewInternalServerError("failed to create book")
	}

	return id, nil
}

// Fetch implements domain.BookUseCase.
func (b *BookUsecase) Fetch(ctx context.Context, params domain.RequestQueryParams) ([]domain.Book, int64, error) {
	books, total, err := b.bookRepo.Fetch(ctx, params)

	if err != nil {
		log.Error().Err(err).Str("layer", "usecase").Msg("failed to fetch book")

		return nil, 0, baseErr.NewInternalServerError("failed to fetch book")
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
