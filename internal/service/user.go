package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"fmt"
	"net/http"
)

func CreateOrUpdateUser(user model.User, conditions ...interface{}) (model.CustomError, model.User) {
	newUser := model.User{
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
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
	err3 := db.CreateOrUpdate(&model.User{}, &newUser, conditions...)
	return err3, newUser
}

func GetUser(conditions ...interface{}) (model.CustomError, model.User) {
	var user model.User
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, user
	}
	defer db.Close()
	result := db.DB().Take(&user, conditions...)
	if result.RowsAffected == 0 {
		return model.ErrInternalDatabaseResourceNotFound, user
	}
	return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), user
}

func UserExists(conditions ...interface{}) bool {
	err, user := GetUser(conditions...)
	return err != model.ErrInternalDatabaseResourceNotFound && user.ID != 0
}

func GetMoviesByUser(username string) (model.CustomError, []model.Movie) {
	var movies []model.Movie
	err, user := GetUser( "username = ?", username)
	if err.IsNotNil() {
		return err, nil
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	var ratings []model.UserRating
	result := db.DB().Find(&ratings, "user_id = ?", user.ID)
	for _, rating := range ratings {
		if rating.MovieId != "" {
			err, movie := omdb.FindMovieByID(rating.MovieId)
			if err.IsNotNil() {
				return err, nil
			}
			movie.Ratings = append(movie.Ratings, model.MovieRating{
				Source: model.RatingSourceCineCircle,
				Value:  fmt.Sprintf("%.1f%s", rating.Rating, model.RatingOver),
				Comment: rating.Comment,
				PostedDate: rating.UpdatedAt,
			})
			movies = append(movies, movie)
		}
	}
	return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), movies
}

func SearchUsers(username, fullname, email string) (model.CustomError, []model.User) {
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	var users []model.User
	queryUsername := "%" + username + "%"
	queryFullname := "%" + fullname + "%"
	queryEmail := "%" + email + "%"
	err2 := db.DB().Find(&users, "username LIKE ? AND full_name LIKE ? AND email LIKE ?",
		queryUsername, queryFullname, queryEmail).Error
	if err2 != nil {
		return model.NewCustomError(err2, model.ErrInternalDatabaseQueryFailed.HttpCode(), model.ErrInternalDatabaseQueryFailedCode), users
	}
	return model.NoErr, users
}