package domain

type Request struct {
	Keyword   string
	Page      int64
	PerPage   int64
	Offset    int64
	SortBy    string
	SortOrder string `defailt:"desc"`
	StartDate string
	EndDate   string
	Filters   map[string]interface{}
}
