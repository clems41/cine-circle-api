package circleRepository

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

type SearchForm struct {
	gormUtils.PaginationQuery
	UserId uint
}

type SearchView struct {
	Total   int64
	Circles []model.Circle
}
