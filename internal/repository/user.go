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
	Email 			string 				`gorm:"uniqueIndex"`
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

func (r userRepository) Create(creation userDom.Creation) (result userDom.Result, err error) {
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

func (r userRepository) Update(update userDom.Update) (result userDom.Result, err error) {
	var user User
	err = r.DB.
		Take(&user, "id = ?", update.UserID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}

	err = r.DB.
		Model(&user).
		Updates(User{Email: update.Email, DisplayName: update.DisplayName}).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}

	result = r.toResult(user)
	return
}

func (r userRepository) UpdatePassword(updatePassword userDom.UpdatePassword) (result userDom.Result, err error) {
	var user User
	err = r.DB.
		Take(&user, "id = ?", updatePassword.UserID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}
	user.HashedPassword = updatePassword.NewHashedPassword

	err = r.DB.
		Model(&user).
		Updates(User{HashedPassword: updatePassword.NewHashedPassword}).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}

	result = r.toResult(user)
	return
}

func (r userRepository) Delete(delete userDom.Delete) (err error) {
	var user User
	err = r.DB.
		Take(&user, "id = ?", delete.UserID).
		Error
	if err != nil {
		return typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}
	err = r.DB.
		Delete(&user).
		Error
	return
}

func (r *userRepository) Get(get userDom.Get) (result userDom.Result, err error) {
	user, err := r.getUser(get)
	if err != nil {
		return
	}

	result = r.toResult(user)
	return
}

func (r userRepository) GetHashedPassword(get userDom.Get) (hashedPassword string, err error) {
	user, err := r.getUser(get)
	if err != nil {
		return
	}

	hashedPassword = user.HashedPassword
	return
}

func (r* userRepository) toResult(user User) (result userDom.Result) {
	return userDom.Result{
		UserID:      user.GetID(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}
}

func (r* userRepository) getUser(get userDom.Get) (user User, err error) {
	query := r.DB

	if get.UserID != 0 {
		query = query.Where("id = ?", get.UserID)
	}
	if get.Username != "" {
		query = query.Where("username = ?", get.Username)
	}
	if get.Email != "" {
		query = query.Where("email = ?", get.Email)
	}

	err = query.
		Take(&user).
		Error

	if err != nil {
		err = typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}
	return
}