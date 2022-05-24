package pgRepositories

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"cine-circle-api/internal/repository/postgres/pgModel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository.User = (*userPgRepository)(nil)
var _ PgRepository = (*userPgRepository)(nil)

type userPgRepository struct {
	DB *gorm.DB
}

func NewUser(DB *gorm.DB) *userPgRepository {
	return &userPgRepository{DB: DB}
}

func (repo *userPgRepository) Migrate() (err error) {
	err = repo.DB.
		AutoMigrate(&pgModel.User{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (repo *userPgRepository) GetFromLogin(login string) (user model.User, ok bool, err error) {
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

func (repo *userPgRepository) Get(userId uint) (user model.User, ok bool, err error) {
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

func (repo *userPgRepository) Save(user *model.User) (err error) {
	err = repo.DB.Save(user).Error
	return errors.WithStack(err)
}

func (repo *userPgRepository) Delete(userId uint) (ok bool, err error) {
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

func (repo *userPgRepository) Search(form repository.UserSearchForm) (view repository.UserSearchView, err error) {
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

func (repo *userPgRepository) UsernameAlreadyExists(username string) (exists bool, err error) {
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

func (repo *userPgRepository) EmailAlreadyExists(email string) (exists bool, err error) {
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
