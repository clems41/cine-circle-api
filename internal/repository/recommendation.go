package repository

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository/postgres/pgModel"
	"cine-circle-api/pkg/utils/searchUtils"
)

type RecommendationSearchForm struct {
	searchUtils.PaginationRequest
	UserId  uint
	MovieId uint
	Type    string
}

type RecommendationSearchView struct {
	searchUtils.Page
	Recommendations []pgModel.Recommendation
}

type Recommendation interface {
	Create(recommendation *model.Recommendation) (err error)
	Search(form RecommendationSearchForm) (view RecommendationSearchView, err error)
}
