package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
)

func AddRating(rating model.UserRating, movieId, username string) (model.CustomError, model.UserRating) {
	if !omdb.MovieExists(movieId) {
		return model.ErrInternalDatabaseResourceNotFound, rating
	}
	if rating.Rating > model.RatingBoundMax || rating.Rating < model.RatingBoundMin {
		return model.ErrInternalApiUnprocessableEntity, rating
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, rating
	}
	rating.UserID = user.ID
	rating.MovieId = movieId
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, rating
	}
	defer db.Close()
	err = db.CreateOrUpdate(&model.UserRating{}, &rating, "movie_id = ? AND user_id = ?", rating.MovieId, rating.UserID)
	return err, rating
}
