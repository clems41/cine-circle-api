package repository

import (
	"cine-circle/internal/domain/recommendationDom"
	"gorm.io/gorm"
)

var _ recommendationDom.Repository = (*recommendationRepository)(nil)

type recommendationRepository struct {
	DB *gorm.DB
}

func NewRecommendationRepository(DB *gorm.DB) *recommendationRepository {
	return &recommendationRepository{DB: DB}
}

func (r recommendationRepository) Migrate() {

}