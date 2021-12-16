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
	AddUser(userId uint, circle *model.Circle) (err error)
	DeleteUser(userId uint, circle *model.Circle) (err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) Get(circleId uint) (circle model.Circle, ok bool, err error) {
	err = repo.DB.
		Preload("Users").
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
	err = repo.DB.
		Omit("Users").
		Create(circle).
		Error
	return errors.WithStack(err)
}

func (repo *repository) Save(circle *model.Circle) (err error) {
	err = repo.DB.
		Omit("Users").
		Save(circle).
		Error
	return errors.WithStack(err)
}

func (repo *repository) Delete(circleId uint) (err error) {
	err = repo.DB.
		Exec("DELETE FROM circle_users where circle_id = ?", circleId).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Delete(&model.Circle{}, circleId).
		Error
	return errors.WithStack(err)
}

func (repo *repository) Search(form SearchForm) (view SearchView, err error) {
	query := repo.DB

	err = query.
		Where("id IN (SELECT circle_id FROM circle_users WHERE user_id = ?)", form.UserId).
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

func (repo *repository) AddUser(userId uint, circle *model.Circle) (err error) {
	err = repo.DB.
		Exec("INSERT INTO circle_users(user_id, circle_id) VALUES (?, ?)", userId, circle.ID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Preload("Users").
		Take(circle).
		Error
	return errors.WithStack(err)
}

func (repo *repository) DeleteUser(userId uint, circle *model.Circle) (err error) {
	err = repo.DB.
		Raw("DELETE FROM circle_users where user_id = ? and circle_id = ?", userId, circle.ID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Preload("Users").
		Take(circle).
		Error
	return errors.WithStack(err)
}
