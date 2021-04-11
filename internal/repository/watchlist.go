package repository

import (
	"cine-circle/internal/domain/watchlistDom"
	"gorm.io/gorm"
)

var _ watchlistDom.Repository = (*watchlistRepository)(nil)

type watchlistRepository struct {
	DB *gorm.DB
}

func NewWatchlistRepository(DB *gorm.DB) *watchlistRepository {
	return &watchlistRepository{DB: DB}
}

func (r watchlistRepository) Migrate() {

}