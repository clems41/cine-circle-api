package searchUtils

// FieldsToColumnConversionMap is used to convert all fields known by request sender (FRONT for example) into column name known by database.
type FieldsToColumnConversionMap map[string]string

type Page struct {
	NumberOfItems int `json:"numberOfItems"`
	NumberOfPages int `json:"numberOfPages"`
	PageSize      int `json:"pageSize"`
	CurrentPage   int `json:"currentPage"`
}

type SortingRequest struct {
	Sort []string `json:"sort"`
}

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
