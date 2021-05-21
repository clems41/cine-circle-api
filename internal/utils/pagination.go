package utils

import (
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"math"
	"strconv"
)

type Page struct {
	NumberOfItems int `json:"number_of_items"`
	NumberOfPages int `json:"number_of_pages"`
	PageSize      int `json:"page_size"`
	CurrentPage   int `json:"current_page"`
}

type SortingRequest struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (request PaginationRequest) Offset() int {

	offset := 0

	if request.Page > 1 {
		offset = (request.Page - 1) * request.PageSize
	}

	return offset
}

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

func (request SortingRequest) OrderSQL() (query string) {
	query = request.Field + " "
	if request.Desc {
		query += "desc"
	} else {
		query += "asc"
	}
	return
}

func ExtractPaginationRequest(req *restful.Request) (request PaginationRequest, err error) {

	page := req.QueryParameter("page")
	if page != "" {

		request.Page, err = strconv.Atoi(page)

		if err != nil {
			return request, typedErrors.NewBadRequestErrorf("'page' parameter is supposed to be an integer: %s", err.Error())
		}

	} else {

		request.Page = 1
	}
	pageSize := req.QueryParameter("page_size")
	if pageSize != "" {

		request.PageSize, err = strconv.Atoi(pageSize)

		if err != nil {
			return request, typedErrors.NewBadRequestErrorf("'page_size' parameter is supposed to be an integer: %s", err.Error())
		}

	} else {

		request.PageSize = 20
	}

	return
}


func ExtractSortingRequest(req *restful.Request, defaultField string, defaultDesc bool) (request SortingRequest, err error) {
	field := req.QueryParameter("field")
	if field != "" {
		request.Field = field
	} else {
		request.Field = defaultField
	}
	desc := req.QueryParameter("desc")
	if desc != "" {
		request.Desc = desc == "true"
	} else {
		request.Desc = defaultDesc
	}
	return
}
