package circleRepository

import (
	"cine-circle-api/internal/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	Get(circleId uint) (circle model.Circle, ok bool, err error)
	Create(circle *model.Circle) (err error)
	Save(circle *model.Circle) (err error)
	Delete(circleId uint) (err error)
	Search(repoForm SearchForm) (repoView SearchView, err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) Get(circleId uint) (circle model.Circle, ok bool, err error) {
	err = repo.DB.
		Take(&circle, circleId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return circle, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) Create(circle *model.Circle) (err error) {
	err = repo.DB.Create(circle).Error
	return errors.WithStack(err)
}

func (repo *repository) Save(circle *model.Circle) (err error) {
	err = repo.DB.Save(circle).Error
	return errors.WithStack(err)
}

func (repo *repository) Delete(circleId uint) (err error) {
	err = repo.DB.Delete(&model.Circle{}, circleId).Error
	return errors.WithStack(err)
}

func (repo *repository) Search(form SearchForm) (view SearchView, err error) {
	query := repo.DB
	if form.CircleName != "" {
		query = query.Where("name ilike ?", "%"+form.CircleName+"%")
	}

	err = query.
		Offset(form.Offset()).
		Limit(form.PageSize).
		Find(&view.Circles).
		Limit(-1).
		Offset(-1).
		Count(&view.Total).
		Error

	if err != nil {
		return view, errors.WithStack(err)
	}
	return
}
