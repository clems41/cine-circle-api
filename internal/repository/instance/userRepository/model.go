package userRepository

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

type SearchForm struct {
	gormUtils.PaginationQuery
	gormUtils.SortQuery
	Keyword string
}

type SearchView struct {
	Total int64
	Users []model.User
}
