package watchlistDom

import (
	"gorm.io/gorm"
)

var _ Repository = (*watchlistRepository)(nil)

type watchlistRepository struct {
	DB *gorm.DB
}

func NewWatchlistRepository(DB *gorm.DB) *watchlistRepository {
	return &watchlistRepository{DB: DB}
}

func (r watchlistRepository) Migrate() {

}
