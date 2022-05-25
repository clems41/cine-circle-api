package repository

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/searchUtils"
)

type RecommendationSearchForm struct {
	searchUtils.PaginationRequest
	UserId  uint
	MediaId uint
	Type    string
}

type Recommendation interface {
	Create(recommendation *model.Recommendation) (err error)
	Search(form RecommendationSearchForm) (recommendations []model.Recommendation, total int64, err error)
}
