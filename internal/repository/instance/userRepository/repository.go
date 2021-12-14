package userRepository

import (
	"cine-circle-api/internal/repository/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	GetUserFromLogin(login string) (user model.User, ok bool, err error)
	GetUser(userId uint) (user model.User, ok bool, err error)
	Create(user *model.User) (err error)
	Save(user *model.User) (err error)
	Delete(userId uint) (ok bool, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) GetUserFromLogin(login string) (user model.User, ok bool, err error) {
	err = repo.DB.
		Take(&user, "username = ? OR email = ?", login, login).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) GetUser(userId uint) (user model.User, ok bool, err error) {
	err = repo.DB.
		Take(&user, userId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *repository) Create(user *model.User) (err error) {
	err = repo.DB.Create(user).Error
	return
}

func (repo *repository) Save(user *model.User) (err error) {
	err = repo.DB.Save(user).Error
	return
}

func (repo *repository) Delete(userId uint) (ok bool, err error) {
	// We should first update unique field to make sure they can be used again after deletion
	var user model.User
	err = repo.DB.
		Take(&user, userId).
		Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	user.Username = uuid.New().String()
	user.Email = uuid.New().String()
	err = repo.DB.
		Save(user).
		Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	err = repo.DB.Delete(&model.User{}, userId).Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	return true, nil
}

func (repo *repository) Search(form SearchForm) (view SearchView, err error) {
	query := repo.DB
	if form.FirstNameKeyword != "" {
		query = query.Where("first_name ilike ?", "%"+form.FirstNameKeyword+"%")
	}
	if form.LastNameKeyword != "" {
		query = query.Where("last_name ilike ?", "%"+form.LastNameKeyword+"%")
	}
	if form.EmailKeyword != "" {
		query = query.Where("email ilike ?", "%"+form.EmailKeyword+"%")
	}
	if form.UsernameKeyword != "" {
		query = query.Where("username ilike ?", "%"+form.UsernameKeyword+"%")
	}
	if form.RoleKeyword != "" {
		query = query.Where("role ilike ?", "%"+form.RoleKeyword+"%")
	}
	if form.ActiveKeyword != "" {
		if form.ActiveKeyword == "true" {
			query = query.Where("active = true")
		} else if form.ActiveKeyword == "false" {
			query = query.Where("active = false")
		}
	}

	err = query.
		Offset(form.Offset()).
		Limit(form.PageSize).
		Order(form.OrderSQL()).
		Find(&view.Users).
		Limit(-1).
		Offset(-1).
		Count(&view.Total).
		Error

	if err != nil {
		return view, errors.WithStack(err)
	}
	return
}
