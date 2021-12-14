package userRepository

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

type SearchForm struct {
	gormUtils.PaginationQuery
	gormUtils.SortQuery
	FirstNameKeyword string
	LastNameKeyword  string
	EmailKeyword     string
	UsernameKeyword  string
	RoleKeyword      string
	ActiveKeyword    string
}

type SearchView struct {
	Total int64
	Users []model.User
}
