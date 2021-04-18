package repository

import (
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
	"gorm.io/gorm"
	"strings"
)

var _ userDom.Repository = (*userRepository)(nil)

type User struct {
	Metadata
	Username 		string 				`gorm:"uniqueIndex"`
	DisplayName 	string
	Email 			string 				`gorm:"index"`
	HashedPassword 	string
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *userRepository {
	return &userRepository{DB: DB}
}

func (r userRepository) Migrate() {

	r.DB.AutoMigrate(&User{})

}

func (r userRepository) CreateUser(creation userDom.Creation) (result userDom.Result, err error) {
	user := User{
		Username:       strings.ToLower(creation.Username),
		DisplayName:    creation.DisplayName,
		Email:          creation.Email,
		HashedPassword: creation.Password,
	}

	err = r.DB.
		Create(&user).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}

	result = userDom.Result{
		UserID:      user.GetID(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}
	return
}

func (r *userRepository) GetUser(get userDom.Get) (result userDom.Result, err error) {
	var user User
	if get.UserID != 0 {
		err = r.DB.
			Find(&user, "id = ?", get.UserID).
			Error
	} else if get.Username != "" {
		err = r.DB.
			Find(&user, "username = ?", get.Username).
			Error
	} else if get.Email != "" {
		err = r.DB.
			Find(&user, "email = ?", get.Email).
			Error
	}
	if err != nil || user.GetID() == 0 {
		return
	}
	result = userDom.Result{
		UserID:      user.GetID(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}
	return
}