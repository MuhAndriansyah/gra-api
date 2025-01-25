package repository

import (
	"backend-layout/internal/domain"
	"fmt"
	"strings"
)

const (
	bookColumns = `
        books.id,
        books.title,
        books.slug,
				books.author_id,
        authors.name as author_name,
				books.publisher_id,
        publishers.name as publisher_name,
        books.publish_year,
        books.total_page,
        books.description,
        books.sku,
        books.isbn,
        books.price,
        books.created_at,
        books.updated_at
    `

	baseQuery = `
        SELECT %s
        FROM books
        JOIN authors ON authors.id = books.author_id
        JOIN publishers ON publishers.id = books.publisher_id
        WHERE 1=1
    `

	countQuery = `
        SELECT COUNT(1) as total_book
        FROM books
        JOIN authors ON authors.id = books.author_id
        JOIN publishers ON publishers.id = books.publisher_id
        WHERE 1=1
    `
)

func buildBookQuery(params domain.RequestQueryParams) (string, []interface{}) {
	var (
		query      = fmt.Sprintf(baseQuery, bookColumns)
		conditions = make([]string, 0)
		args       = make([]interface{}, 0)
		argCounter = 1
	)

	// Handle price filters
	if v, ok := params.Filters["min_price"]; ok && v.(int64) > 0 {
		conditions = append(conditions, fmt.Sprintf("books.price >= $%d", argCounter))
		args = append(args, v.(int64))
		argCounter++
	}

	if v, ok := params.Filters["max_price"]; ok && v.(int64) > 0 {
		conditions = append(conditions, fmt.Sprintf("books.price <= $%d", argCounter))
		args = append(args, v.(int64))
		argCounter++
	}

	// Handle search
	if params.Keyword != "" {
		conditions = append(conditions, fmt.Sprintf("books.title ILIKE $%d", argCounter))
		args = append(args, "%"+escapeSQLLike(params.Keyword)+"%")
		argCounter++
	}

	// Append conditions jika ada
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Handle sorting
	query += buildOrderBy(params.SortBy)

	// Handle pagination
	offset := (params.Page - 1) * params.PerPage
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, params.PerPage, offset)

	return query, args
}

func buildOrderBy(sortBy string) string {
	switch sortBy {
	case "highest_price":
		return " ORDER BY books.price DESC"
	case "lowest_price":
		return " ORDER BY books.price ASC"
	default:
		return " ORDER BY books.created_at DESC"
	}
}

func buildCountBookQuery(params domain.RequestQueryParams) (string, []interface{}) {
	var (
		query      = countQuery
		conditions = make([]string, 0)
		args       = make([]interface{}, 0)
		argCounter = 1
	)

	// Handle price filters
	if v, ok := params.Filters["min_price"]; ok && v.(int64) > 0 {
		conditions = append(conditions, fmt.Sprintf("books.price >= $%d", argCounter))
		args = append(args, v.(int64))
		argCounter++
	}

	if v, ok := params.Filters["max_price"]; ok && v.(int64) > 0 {
		conditions = append(conditions, fmt.Sprintf("books.price <= $%d", argCounter))
		args = append(args, v.(int64))
		argCounter++
	}

	// Handle search
	if params.Keyword != "" {
		conditions = append(conditions, fmt.Sprintf("books.title ILIKE $%d", argCounter))
		args = append(args, "%"+escapeSQLLike(params.Keyword)+"%")
		argCounter++
	}

	// Append conditions jika ada
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query, args
}

func escapeSQLLike(input string) string {
	return strings.NewReplacer(
		"%", "\\%",
		"_", "\\_",
	).Replace(input)
}
