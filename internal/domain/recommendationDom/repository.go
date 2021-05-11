package recommendationDom

import (
	"gorm.io/gorm"
)

var _ Repository = (*recommendationRepository)(nil)

type recommendationRepository struct {
	DB *gorm.DB
}

func NewRecommendationRepository(DB *gorm.DB) *recommendationRepository {
	return &recommendationRepository{DB: DB}
}

func (r recommendationRepository) Migrate() {

}
