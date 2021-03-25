package service

import (
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"gorm.io/gorm"
	"net/http"
)

func CreateOrUpdateCircle(circle model.Circle) (model.CustomError, model.Circle) {
	newCircle := model.Circle{
		Model:       gorm.Model{ID: circle.ID},
		Name:        circle.Name,
		Description: circle.Description,
	}
	err := newCircle.IsValid()
	if err.IsNotNil() {
		return err, newCircle
	}
	db, err2 := database.OpenConnection()
	if err2.IsNotNil() {
		return err2, newCircle
	}
	defer db.Close()
	err3 := db.CreateOrUpdate(&model.Circle{}, &newCircle, "id = ?", newCircle.ID)
	return err3, newCircle
}

func DeleteCircle(circleId uint) model.CustomError {
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	result := db.DB().Delete(&model.Circle{}, "id = ?", circleId)
	if result.Error != nil {
		return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode)
	}
	if result.RowsAffected != 1 {
		return model.ErrInternalDatabaseResourceNotFound
	}
	return model.NoErr
}

func AddUserToCircle(circleId, userId uint) (model.CustomError, model.Circle) {
	var circle model.Circle
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, circle
	}
	defer db.Close()
	db.DB().Exec("INSERT INTO user_circle(user_id, circle_id) VALUES (?, ?)", userId, circleId)
	db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	return model.NoErr, circle
}

func RemoveUserFromCircle(circleId, userId uint) (model.CustomError, model.Circle) {
	var circle model.Circle
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, circle
	}
	defer db.Close()
	db.DB().Exec("DELETE FROM user_circle WHERE user_id = ? AND circle_id = ?", userId, circleId)
	db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	return model.NoErr, circle
}

func GetCircles(username string, conditions ...interface{}) (model.CustomError, []model.Circle) {
	var circles []model.Circle
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	err2 := db.DB().Preload("Users").Find(&circles, conditions...).Error
	if err2 != nil {
		return model.NewCustomError(err2, model.ErrInternalDatabaseConnectionFailed.HttpCode(), model.ErrInternalDatabaseConnectionFailedCode), nil
	}
	var result []model.Circle
	for _, circle := range circles {
		for _, user := range circle.Users {
			if user.Username == username {
				result = append(result, circle)
				break
			}
		}
	}
	return model.NoErr, result
}

func GetMoviesForCircle(circleId uint, sort string) (model.CustomError, []model.Movie) {
	var circle model.Circle
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	result := db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	if result.RowsAffected != 1 {
		return model.ErrInternalDatabaseResourceNotFound, nil
	}
	var movies []model.Movie
	for _, user := range circle.Users {
		err2, userMovies := GetMoviesByUser("username = ?", user.Username)
		if err2.IsNotNil() {
			return err2, nil
		}
		movies = append(movies, userMovies...)
	}
	err2, moviesMerged := MergeMovies(movies)
	if err2.IsNotNil() {
		return err2, nil
	}
	err2 = SortMovies(moviesMerged, sort)
	return err2, moviesMerged
}