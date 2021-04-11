package repository

import (
	"cine-circle/internal/domain/movieDom"
	"gorm.io/gorm"
)

var _ movieDom.Repository = (*movieRepository)(nil)

type movieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository(DB *gorm.DB) *movieRepository {
	return &movieRepository{DB: DB}
}

func (r movieRepository) Migrate() {

}