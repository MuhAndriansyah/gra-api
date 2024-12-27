package helper

import (
	"backend-layout/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetRequestParams(c echo.Context) domain.Request {
	sortOrder := c.QueryParam("sort_order")
	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.QueryParam("per_page"), 10, 64)

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	params := domain.Request{
		Keyword:   c.QueryParam("q"),
		Page:      page,
		PerPage:   perPage,
		Offset:    offset,
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: sortOrder,
		StartDate: c.QueryParam("start_date"),
		EndDate:   c.QueryParam("end_date"),
	}

	return params
}
