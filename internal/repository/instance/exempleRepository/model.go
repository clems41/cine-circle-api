package exempleRepository

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

type SearchForm struct {
	gormUtils.PaginationQuery
	gormUtils.SortQuery
	// TODO add your keyword fields here (cf. userRepository example)
}

type SearchView struct {
	Total    int64
	Exemples []model.Exemple
}
