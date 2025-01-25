package helper

import (
	"backend-layout/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetRequestParams(c echo.Context) domain.RequestQueryParams {
	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	perPage, _ := strconv.ParseInt(c.QueryParam("per_page"), 10, 64)

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 10
	}

	params := domain.RequestQueryParams{
		Keyword:   c.QueryParam("q"),
		Page:      page,
		PerPage:   perPage,
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
		StartDate: c.QueryParam("start_date"),
		EndDate:   c.QueryParam("end_date"),
	}

	return params
}
