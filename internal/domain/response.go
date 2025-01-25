package domain

type ResponseBody struct {
	Data interface{} `json:"data,omitempty"`
	Meta *Pagination `json:"meta,omitempty"`
}

type Pagination struct {
	TotalCount  int64   `json:"total_count"`
	TotalPage   float64 `json:"total_page"`
	CurrentPage int64   `json:"current_page"`
	PerPage     int64   `json:"per_page"`
}
