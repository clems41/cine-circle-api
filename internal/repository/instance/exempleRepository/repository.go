package exempleRepository

import (
	"cine-circle-api/internal/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	Get(exempleId uint) (exemple model.Exemple, ok bool, err error)
	Create(exemple *model.Exemple) (err error)
	Save(exemple *model.Exemple) (err error)
	Delete(exempleId uint) (err error)
	Search(repoForm SearchForm) (repoView SearchView, err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) Get(exempleId uint) (exemple model.Exemple, ok bool, err error) {
	err = repo.DB.
		Take(&exemple, exempleId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exemple, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) Create(exemple *model.Exemple) (err error) {
	err = repo.DB.
		Create(exemple).
		Error
	return
}

func (repo *repository) Save(exemple *model.Exemple) (err error) {
	err = repo.DB.
		Save(exemple).
		Error
	return
}

func (repo *repository) Delete(exempleId uint) (err error) {
	// We should first update unique field to make sure they can be used again after deletion
	var exemple model.Exemple
	err = repo.DB.
		Take(&exemple, exempleId).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Delete(&model.Exemple{}, exempleId).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (repo *repository) Search(form SearchForm) (view SearchView, err error) {
	query := repo.DB

	err = query.
		Offset(form.Offset()).
		Limit(form.PageSize).
		Order(form.OrderSQL()).
		Find(&view.Exemples).
		Limit(-1).
		Offset(-1).
		Count(&view.Total).
		Error

	if err != nil {
		return view, errors.WithStack(err)
	}
	return
}
