package service

import (
	"cine-circle/internal/model"
	"cine-circle/internal/repository"
	"cine-circle/internal/typedErrors"
)

func CreateCircle(circle model.Circle, username string) (typedErrors.CustomError, model.Circle) {
	newCircle := model.Circle{
		Name:        circle.Name,
		Description: circle.Description,
	}
	err := newCircle.IsValid()
	if err.IsNotNil() {
		return err, newCircle
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, newCircle
	}
	newCircle.Users = []model.User{user}
	db, err2 := repository.OpenConnection()
	if err2.IsNotNil() {
		return err2, newCircle
	}
	defer db.Close()
	err3 := db.DB().Create(&newCircle).Association("Users").Error
	return typedErrors.NewRepositoryQueryFailedError(err3), newCircle
}

func UpdateCircle(circle model.Circle, circleId, username string) (typedErrors.CustomError, model.Circle) {
	err := circle.IsValid()
	if err.IsNotNil() {
		return err, circle
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, circle
	}
	db, err2 := repository.OpenConnection()
	if err2.IsNotNil() {
		return err2, circle
	}
	defer db.Close()
	var association model.UserCircle
	result := db.DB().Table("user_circle").Find(&association, "circle_id = ? AND user_id = ?", circleId, user.ID)
	if result.RowsAffected != 1 {
		return typedErrors.ErrApiUserBadCredentials, circle
	}
	err3 := db.DB().Model(&circle).Where("id = ?", circleId).Update("name", circle.Name).Update("description", circle.Description).Error
	if err3 == nil {
		err3 = db.DB().Preload("Users").Find(&circle, "id = ?", circleId).Error
	}
	return typedErrors.NewRepositoryQueryFailedError(err3), circle
}

func DeleteCircle(circleId, username string) typedErrors.CustomError {
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err
	}
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	var association model.UserCircle
	result := db.DB().Table("user_circle").Find(&association, "circle_id = ? AND user_id = ?", circleId, user.ID)
	if result.RowsAffected != 1 {
		return typedErrors.ErrApiUserBadCredentials
	}
	result = db.DB().Delete(&model.Circle{}, "id = ?", circleId)
	if result.Error != nil {
		return typedErrors.NewRepositoryQueryFailedError(result.Error)
	}
	if result.RowsAffected != 1 {
		return typedErrors.ErrRepositoryResourceNotFound
	}
	result = db.DB().Table("user_circle").Delete(&model.UserCircle{}, "circle_id = ?", circleId)
	if result.Error != nil {
		return typedErrors.NewRepositoryQueryFailedError(result.Error)
	}
	return typedErrors.NoErr
}

func AddUserToCircle(circleId, userId uint) (typedErrors.CustomError, model.Circle) {
	var circle model.Circle
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, circle
	}
	defer db.Close()
	db.DB().Exec("INSERT INTO user_circle(user_id, circle_id) VALUES (?, ?)", userId, circleId)
	db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	return typedErrors.NoErr, circle
}

func RemoveUserFromCircle(circleId, userId uint) (typedErrors.CustomError, model.Circle) {
	var circle model.Circle
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, circle
	}
	defer db.Close()
	db.DB().Exec("DELETE FROM user_circle WHERE user_id = ? AND circle_id = ?", userId, circleId)
	db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	return typedErrors.NoErr, circle
}

func GetCircles(username string, conditions ...interface{}) (typedErrors.CustomError, []model.Circle) {
	var circles []model.Circle
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	err2 := db.DB().Preload("Users").Find(&circles, conditions...).Error
	if err2 != nil {
		return typedErrors.NewRepositoryQueryFailedError(err2), nil
	}
	result := []model.Circle{}
	for _, circle := range circles {
		for _, user := range circle.Users {
			if user.Username == username {
				result = append(result, circle)
				break
			}
		}
	}
	return typedErrors.NoErr, result
}

func GetMoviesForCircle(circleId uint, sort string) (typedErrors.CustomError, []model.Movie) {
	var circle model.Circle
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, nil
	}
	defer db.Close()
	result := db.DB().Preload("Users").Take(&circle, "id = ?", circleId)
	if result.RowsAffected != 1 {
		return typedErrors.ErrRepositoryResourceNotFound, nil
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