package mediaRepository

import (
	"cine-circle-api/internal/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	GetMovie(movieId uint) (movie model.Movie, ok bool, err error)
	Save(movie *model.Movie) (err error)
	Create(movie *model.Movie) (err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) GetMovie(movieId uint) (movie model.Movie, ok bool, err error) {
	err = repo.DB.
		Take(&movie, movieId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return movie, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) Save(movie *model.Movie) (err error) {
	err = repo.DB.Save(movie).Error
	return
}

func (repo *repository) Create(movie *model.Movie) (err error) {
	err = repo.DB.Create(movie).Error
	return
}
