package helper

import (
	"backend-layout/internal/domain"
	"math"

	"github.com/labstack/echo/v4"
)

func Paginate(_ echo.Context, data interface{}, total int64, params domain.Request) *domain.ResponseBody {
	return &domain.ResponseBody{
		Data: data,
		Meta: &domain.Pagination{
			TotalCount:  total,
			TotalPage:   math.Ceil(float64(total) / float64(params.PerPage)),
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
		},
	}
}
