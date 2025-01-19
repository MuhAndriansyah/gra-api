package repository

import (
	"backend-layout/internal/domain"
	"fmt"
	"strings"
)

const (
	baseQuery = `SELECT id, title, slug, author.name as author_name, publisher.name as publisher_name, publish_year, total_page, description, sku, isbn, discount, price, created_at, updated_at 
				 FROM books 
				 JOIN authors ON authors.id = books.author_id 
				 JOIN publishers ON publishers.id = books.publisher_id
				 WHERE 1=1`
)

func buildBookQuery(params domain.BookQueryParams) (string, []interface{}) {

	var queryBuilder strings.Builder
	queryBuilder.WriteString(baseQuery)

	args := make([]interface{}, 0)
	argCounter := 1

	if params.MinPrice > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND books.price >= $%d", argCounter))
		args = append(args, params.MinPrice)
		argCounter++
	}

	if params.MaxPrice > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND price <= $%d", argCounter))
		args = append(args, params.MaxPrice)
		argCounter++
	}

	if params.Search != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND books.title ILIKE $%d", argCounter))
		args = append(args, "%"+escapeSQLLike(params.Search)+"%")
		argCounter++
	}

	if params.Sort == "highest_price" {
		queryBuilder.WriteString(" ORDER BY books.price DESC")
	} else if params.Sort == "lowest_price" {
		queryBuilder.WriteString(" ORDER BY books.price ASC")
	}

	queryBuilder.WriteString(" ORDER BY books.created_at DESC")

	if params.Page > 0 && params.PerPage > 0 {
		offset := (params.Page - 1) * params.PerPage
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1))
		args = append(args, params.PerPage, offset)
		argCounter += 2
	} else if params.Offset >= 0 && params.PerPage > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1))
		args = append(args, params.PerPage, params.Offset)
		argCounter += 2
	}

	return queryBuilder.String(), args
}

func escapeSQLLike(input string) string {
	replacements := map[string]string{
		"%": "\\%",
		"_": "\\_",
	}
	for old, new := range replacements {
		input = strings.ReplaceAll(input, old, new)
	}
	return input
}
