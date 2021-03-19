package service

import (
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"net/http"
)

func CreateUser(fullName string) model.CustomError {
	user := model.User{FullName: fullName}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	err2 := db.DB().Create(&user).Error
	return model.NewCustomError(err2, http.StatusBadRequest, model.ErrInternalDatabaseCreationFailedCode)
}

func GetUser(id string) (model.CustomError, model.User) {
	var user model.User
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, user
	}
	defer db.Close()
	result := db.DB().First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		return model.ErrInternalApiNotFound, user
	}
	return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), user
}

