package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/model"
	"cine-circle/internal/repository"
	"cine-circle/internal/typedErrors"
)

func CreateOrUpdateUser(user model.User, conditions ...interface{}) (typedErrors.CustomError, model.User) {
	err := user.IsValid()
	if err.IsNotNil() {
		return err, user
	}
	newUser := model.User{
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}
	err = HashAndSaltPassword(user.Password, &newUser)
	if err.IsNotNil() {
		return err, user
	}
	db, err2 := repository.OpenConnection()
	if err2.IsNotNil() {
		return err2, user
	}
	defer db.Close()
	err3 := db.CreateOrUpdate(&model.User{}, &newUser, conditions...)
	return err3, newUser
}

func DeleteUser(conditions ...interface{}) typedErrors.CustomError {
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	err2 := db.DB().Delete(&model.User{}, conditions...).Error
	return typedErrors.NewRepositoryQueryFailedError(err2)
}

func GetUser(conditions ...interface{}) (typedErrors.CustomError, model.User) {
	var user model.User
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, user
	}
	defer db.Close()
	result := db.DB().Take(&user, conditions...)
	if result.RowsAffected == 0 {
		return typedErrors.ErrRepositoryResourceNotFound, user
	}
	return typedErrors.NewRepositoryQueryFailedError(result.Error), user
}

func UserExists(conditions ...interface{}) bool {
	err, user := GetUser(conditions...)
	return err != typedErrors.ErrRepositoryResourceNotFound && user.ID != 0
}

func GetMoviesByUser(conditions ...interface{}) (typedErrors.CustomError, []model.Movie) {
	movies := []model.Movie{}
	err, user := GetUser(conditions...)
	if err.IsNotNil() {
		return err, nil
	}
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	var ratings []model.Rating
	result := db.DB().Find(&ratings, "user_id = ?", user.ID)
	for _, rating := range ratings {
		if rating.MovieID != "" {
			err, movie := omdb.FindMovieByID(rating.MovieID)
			if err.IsNotNil() {
				return err, nil
			}
			movie.UserRatings = append(movie.UserRatings, rating)
			movies = append(movies, movie)
		}
	}
	return typedErrors.NewRepositoryQueryFailedError(result.Error), movies
}

func SearchUsers(username, fullname, email string) (typedErrors.CustomError, []model.User) {
	db, err := repository.OpenConnection()
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
		return typedErrors.NewRepositoryQueryFailedError(err2), users
	}
	return typedErrors.NoErr, users
}