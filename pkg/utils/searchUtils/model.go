package searchUtils

type Page struct {
	NumberOfItems int64 `json:"numberOfItems"`
	NumberOfPages int   `json:"numberOfPages"`
	PageSize      int   `json:"pageSize"`
	CurrentPage   int   `json:"currentPage"`
}

type SortingRequest struct {
	Sort []string `json:"sort"`
}

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
