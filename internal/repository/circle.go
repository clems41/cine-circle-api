package repository

import (
	"cine-circle/internal/domain/circleDom"
	"gorm.io/gorm"
)

var _ circleDom.Repository = (*circleRepository)(nil)

type circleRepository struct {
	DB *gorm.DB
}

func NewCircleRepository(DB *gorm.DB) *circleRepository {
	return &circleRepository{DB: DB}
}

func (r circleRepository) Migrate() {

}