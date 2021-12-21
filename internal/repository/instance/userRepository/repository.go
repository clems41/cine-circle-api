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
	UsernameAlreadyExists(username string) (exists bool, err error)
	EmailAlreadyExists(email string) (exists bool, err error)
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
	return errors.WithStack(err)
}

func (repo *repository) Save(user *model.User) (err error) {
	err = repo.DB.Save(user).Error
	return errors.WithStack(err)
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
	if form.Keyword != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR username ilike ?",
			"%"+form.Keyword+"%", "%"+form.Keyword+"%", "%"+form.Keyword+"%")
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

func (repo *repository) UsernameAlreadyExists(username string) (exists bool, err error) {
	err = repo.DB.
		Select("username").
		Take(&model.User{}, "username = ?", username).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, errors.WithStack(err)
		}
	}
	return true, nil
}

func (repo *repository) EmailAlreadyExists(email string) (exists bool, err error) {
	err = repo.DB.
		Select("email").
		Take(&model.User{}, "email = ?", email).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, errors.WithStack(err)
		}
	}
	return true, nil
}
