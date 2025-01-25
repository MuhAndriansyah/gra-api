package domain

type RequestQueryParams struct {
	Keyword   string
	Page      int64
	PerPage   int64
	SortBy    string
	SortOrder string `defailt:"desc"`
	StartDate string
	EndDate   string
	Filters   map[string]interface{}
}
