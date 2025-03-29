package http

import (
	"backend-layout/helper"
	"backend-layout/internal/domain"
	"backend-layout/internal/middleware"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type BookHandler struct {
	bookUsecase domain.BookUsecase
	rbac        *middleware.RBACMiddleware
}

func NewBookHandler(p *echo.Group, r *echo.Group, bu domain.BookUsecase, rbac *middleware.RBACMiddleware) {
	handler := &BookHandler{
		bookUsecase: bu,
		rbac:        rbac,
	}

	p.GET("/books", handler.List)
	r.POST("/books", handler.Store, rbac.RequiredPermission("book:create"))
	p.GET("/books/:id", handler.Get)
	r.DELETE("/books/:id", handler.Delete, rbac.RequiredPermission("book:delete"))
	r.PATCH("/books/:id", handler.Update, rbac.RequiredPermission("book:update"))
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
		log.Err(err).Msg("failed to fetch book")

		return err
	}

	var booksResponse = make([]domain.BookResponse, len(listBooks))

	for i, v := range listBooks {
		booksResponse[i] = domain.BookToResponse(&v)
	}

	res := helper.Paginate(c, booksResponse, total, params)

	return c.JSON(http.StatusOK, res)
}

func (h *BookHandler) Store(c echo.Context) error {
	req := new(domain.StoreBookRequest)

	correlationID, ok := c.Request().Context().Value(middleware.CorrelationIDKey).(string)
	if !ok {
		correlationID = "unknown"
	}

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	id, err := h.bookUsecase.Store(ctx, req)

	if err != nil {

		log.Err(err).
			Str("coreelation_id", correlationID).
			Str("service", "book-service").
			Str("layer", "handler").
			Str("func", "store").
			Str("method", "POST").
			Str("path", c.Path()).
			Msg("failed to create book")

		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Book created succesfully", "data": echo.Map{
		"id": id,
	}})
}

func (h *BookHandler) Get(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID format")
	}

	ctx := c.Request().Context()

	book, err := h.bookUsecase.Get(ctx, id)

	if err != nil {
		log.Err(err).Msg("failed to get book")
		return err
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Delete(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID format")
	}

	ctx := c.Request().Context()

	err = h.bookUsecase.Delete(ctx, id)

	if err != nil {
		log.Err(err).Msg("failed to delete book")
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "book deleted succesfully"})

}

func (h *BookHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID format")
	}

	req := new(domain.UpdateBookRequest)
	req.ID = id

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	err = h.bookUsecase.Update(ctx, req)

	if err != nil {
		log.Err(err).Msg("failed to update book")
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "book updated successfully"})
}
