package recommendationRepository

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

type SearchForm struct {
	gormUtils.PaginationQuery
	UserId  uint
	MovieId uint
	Type    string
}

type SearchView struct {
	Total           int64
	Recommendations []model.Recommendation
}
