package watchlistDom

import (
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {
}
