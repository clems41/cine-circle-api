package searchUtils

import (
	"cine-circle-api/pkg/utils/stringUtils"
	"fmt"
	"math"
	"strings"
)

// BuildResult returns page object with right fields based on numberOfItems
func (request PaginationRequest) BuildResult(numberOfItems int64) (page Page) {

	page.CurrentPage = request.Page

	if page.CurrentPage < 1 {
		page.CurrentPage = 1
	}

	page.NumberOfItems = numberOfItems

	if request.PageSize < 1 {

		page.PageSize = int(page.NumberOfItems)

		page.NumberOfPages = 1

	} else {

		page.PageSize = request.PageSize

		page.NumberOfPages = int(math.Ceil(float64(page.NumberOfItems) / float64(page.PageSize)))
	}

	return
}

// Offset return offset to start getting result based on request.PaginationRequest. Can be used as SQL offset.
func (request PaginationRequest) Offset() int {

	offset := 0

	if request.Page > 1 {
		offset = (request.Page - 1) * request.PageSize
	}

	return offset
}

// OrderSQL return string that can be used as Order condition for SQL query.
// Sort is a string array which each string must be like 'name:asc' (insensitive cases) to be considered.
// If there is no ':', string will be taken as field name and order will be ASC by default.
func (request SortingRequest) OrderSQL() string {
	var orders []string
	for _, sort := range request.Sort {
		split := strings.Split(sort, ":")
		var ascString, fieldName string
		if len(split) >= 2 {
			if strings.ToLower(split[1]) == "asc" {
				ascString = "asc"
			} else {
				ascString = "desc"
			}
			fieldName = stringUtils.ToSnakeCase(split[0])
		} else {
			ascString = "asc"
			fieldName = stringUtils.ToSnakeCase(sort)
		}
		orders = append(orders, fmt.Sprintf("%s %s", fieldName, ascString))
	}
	return strings.Join(orders, ",")
}
