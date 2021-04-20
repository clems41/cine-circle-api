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
	user.Email = update.Email
	user.DisplayName = update.DisplayName

	err = r.DB.
		Save(&user).
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
		Save(&user).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}

	result = r.toResult(user)
	return
}

func (r userRepository) Delete(delete userDom.Delete) (err error) {
	err = r.DB.
		Delete(&User{}, "id = ?", delete.UserID).
		Error
	if err != nil {
		return typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
	}
	return
}

func (r *userRepository) Get(get userDom.Get) (result userDom.Result, err error) {
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
	result = r.toResult(user)
	return
}

func (r userRepository) GetHashedPassword(username string) (hashedPassword string, err error) {
	var user User
	err = r.DB.
		Select("hashed_password").
		Find(&user, "username = ?", username).
		Error
	if err != nil {
		return hashedPassword, typedErrors.NewRepositoryQueryFailedErrorf(err.Error())
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