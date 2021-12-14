package mediaRepository

import (
	"cine-circle-api/internal/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	Get(mediaId uint) (media model.Media, ok bool, err error)
	Save(media *model.Media) (err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) Get(mediaId uint) (media model.Media, ok bool, err error) {
	err = repo.DB.
		Take(&media, mediaId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return media, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) Save(media *model.Media) (err error) {
	err = repo.DB.Save(media).Error
	return
}
