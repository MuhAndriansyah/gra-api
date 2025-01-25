package http

import (
	"backend-layout/helper"
	"backend-layout/internal/domain"
	"backend-layout/internal/module/book/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type BookHandler struct {
	bookUsecase domain.BookUsecase
	logger      zerolog.Logger
}

func NewUserHanlder(p *echo.Group, r *echo.Group, bu domain.BookUsecase, logger zerolog.Logger) {
	handler := &BookHandler{
		bookUsecase: bu,
		logger:      logger,
	}

	p.GET("/books", handler.List)
	p.POST("/books", handler.Store)
}

func (h *BookHandler) List(c echo.Context) (err error) {
	ctx := c.Request().Context()

	params := helper.GetRequestParams(c)
	minPrice, _ := strconv.ParseInt(c.QueryParam("min_price"), 10, 64)
	maxPrice, _ := strconv.ParseInt(c.QueryParam("max_price"), 10, 64)

	params.Filters = map[string]interface{}{
		"min_price": minPrice,
		"max_price": maxPrice,
	}

	listBooks, total, err := h.bookUsecase.Fetch(ctx, params)

	if err != nil {
		h.logger.Err(err).
			Str("book_usecase", "fetch").
			Str("book_handler", "list").
			Msg("failed to create book")

		return err
	}

	var listBookRes = make([]domain.BookResponse, len(listBooks))

	for i, v := range listBooks {
		listBookRes[i] = *domain.BookToResponse(&v)
	}

	res := helper.Paginate(c, listBookRes, total, params)

	return c.JSON(http.StatusOK, res)
}

func (h *BookHandler) Store(c echo.Context) error {
	req := new(domain.StoreBookRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	id, err := h.bookUsecase.Store(ctx, req)

	if err != nil {
		h.logger.Err(err).
			Str("book_usecase", "store").
			Str("handler", "store").
			Msg("failed to create book")

		if err == repository.ErrISBNDuplicateEntry {
			return echo.NewHTTPError(http.StatusConflict, "isbn already exists")
		}

		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Book created succesfully", "data": echo.Map{
		"id": id,
	}})
}
