package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"fmt"
)

const (
	ratingBoundMin = 0
	ratingBoundMax = 10
)

var (
	ratingOver = "/" + fmt.Sprintf("%d", ratingBoundMax)
)

func AddRating(rating model.UserRating, movieId, username string) (model.CustomError, model.UserRating) {
	if !omdb.MovieExists(movieId) {
		return model.ErrInternalDatabaseResourceNotFound, rating
	}
	if rating.Rating > ratingBoundMax || rating.Rating < ratingBoundMin {
		return model.ErrInternalApiUnprocessableEntity, rating
	}
	err, userId := GetUserIdByUsername(username)
	if err.IsNotNil() {
		return err, rating
	}
	rating.UserID = userId
	rating.MovieId = movieId
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, rating
	}
	defer db.Close()
	err = db.CreateOrUpdate(&model.UserRating{}, &rating, "movie_id = ? AND user_id = ?", rating.MovieId, rating.UserID)
	return err, rating
}
