package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
)

func AddRating(rating model.Rating, movieId, username string) (model.CustomError, model.Rating) {
	if !omdb.MovieExists(movieId) {
		return model.ErrInternalDatabaseResourceNotFound, rating
	}
	if rating.Value > model.RatingBoundMax || rating.Value < model.RatingBoundMin {
		return model.ErrInternalApiUnprocessableEntity, rating
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, rating
	}
	rating.UserID = user.ID
	rating.MovieID = movieId
	rating.Source = model.CineCircleSource
	rating.Username = user.Username
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, rating
	}
	defer db.Close()
	err = db.CreateOrUpdate(&model.Rating{}, &rating, "movie_id = ? AND user_id = ?", rating.MovieID, rating.UserID)
	return err, rating
}

func AddUserRating(username string, movie *model.Movie) model.CustomError {
	var rating model.Rating
	err, user := GetUser("username = ?", username)
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	result := db.DB().Find(&rating, "user_id = ? AND movie_id = ?", user.ID, movie.ID)
	if result.RowsAffected > 0 {
		movie.UserRatings = append(movie.UserRatings, rating)
	}
	return err

}