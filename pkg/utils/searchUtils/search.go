package searchUtils

import (
	"math"
)

// BuildResult returns page object with right fields based on numberOfItems
func (request PaginationRequest) BuildResult(numberOfItems int64) (page Page) {

	page.CurrentPage = request.Page

	if page.CurrentPage < 1 {
		page.CurrentPage = 1
	}

	page.NumberOfItems = int(numberOfItems)

	if request.PageSize < 1 {

		page.PageSize = page.NumberOfItems

		page.NumberOfPages = 1

	} else {

		page.PageSize = request.PageSize

		page.NumberOfPages = int(math.Ceil(float64(page.NumberOfItems) / float64(page.PageSize)))
	}

	return
}
