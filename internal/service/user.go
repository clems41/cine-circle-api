package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"net/http"
)

func CreateUser(username, fullName, email string) (model.CustomError, model.User) {
	user := model.User{
		FullName: fullName,
		Username: username,
		Email:    email,
	}
	err := user.IsValid()
	if err.IsNotNil() {
		return err, user
	}
	db, err2 := database.OpenConnection()
	if err2.IsNotNil() {
		return err2, user
	}
	defer db.Close()
	err3 := db.DB().Create(&user).Error
	return model.NewCustomError(err3, http.StatusBadRequest, model.ErrInternalDatabaseCreationFailedCode), user
}

func GetUser(username string) (model.CustomError, model.User) {
	var user model.User
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, user
	}
	defer db.Close()
	result := db.DB().Take(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		return model.ErrInternalDatabaseResourceNotFound, user
	}
	return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), user
}

func GetUserIdByUsername(username string) (model.CustomError, uint) {
	err, user := GetUser(username)
	return err, user.ID
}

func UsernameAlreadyExists(username string) bool {
	err, user := GetUser(username)
	if err == model.ErrInternalDatabaseResourceNotFound {
		return false
	}
	return user.Username == username
}

func GetMovieByUser(username string) (model.CustomError, []model.Movie) {
	var movies []model.Movie
	err, userId := GetUserIdByUsername(username)
	if err.IsNotNil() {
		return err, nil
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	var ratings []model.UserRating
	result := db.DB().Find(&ratings, "user_id = ?", userId)
	for _, rating := range ratings {
		if rating.MovieId != "" {
			err, movie := omdb.FindMovieByID(rating.MovieId)
			if err.IsNotNil() {
				return err, nil
			}
			movies = append(movies, movie)
		}
	}
	return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), movies
}