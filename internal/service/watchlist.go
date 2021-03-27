package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
)

func AddMovieToWatchlist(username, movieId string) model.CustomError {
	if !omdb.MovieExists(movieId) {
		return model.ErrInternalDatabaseResourceNotFound
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	watchlist := model.Watchlist{
		UserID:  user.ID,
		MovieID: movieId,
	}
	err2 := db.DB().Save(&watchlist).Error
	if err2 != nil {
		return model.NewCustomError(err2, model.ErrInternalDatabaseQueryFailed.HttpCode(), model.ErrInternalDatabaseQueryFailedCode)
	}
	return model.NoErr
}

func RemoveMovieFromWatchlist(username, movieId string) model.CustomError {
	if !omdb.MovieExists(movieId) {
		return model.ErrInternalDatabaseResourceNotFound
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	err2 := db.DB().Delete(&model.Watchlist{}, "user_id = ? AND movie_ID = ?", user.ID, movieId).Error
	if err2 != nil {
		return model.NewCustomError(err2, model.ErrInternalDatabaseQueryFailed.HttpCode(), model.ErrInternalDatabaseQueryFailedCode)
	}
	return model.NoErr
}

func GetMoviesFromWatchlist(username string) (model.CustomError, model.MovieSearch) {
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, model.MovieSearch{}
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, model.MovieSearch{}
	}
	defer db.Close()
	var watchlists []model.Watchlist
	db.DB().Find(&watchlists, "user_id = ?", user.ID)
	var result model.MovieSearch
	for _, watchlist := range watchlists {
		err, movie := omdb.FindMovieByID(watchlist.MovieID)
		if err.IsNotNil() {
			return err, model.MovieSearch{}
		}
		result.Search = append(result.Search, movie.MovieShort())
	}
	if result.Search == nil {
		result.Search = []model.MovieShort{}
	}
	result.TotalResults = len(result.Search)
	return model.NoErr, result
}

func IsInWatchlist(username, movieId string) (model.CustomError, bool) {
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, false
	}
	db, err := database.OpenConnection()
	if err.IsNotNil() {
		return err, false
	}
	defer db.Close()
	var watchlists []model.Watchlist
	result := db.DB().Find(&watchlists, "user_id = ? AND movie_id = ?", user.ID, movieId)
	if result.Error != nil {
		return model.NewCustomError(result.Error, model.ErrInternalDatabaseQueryFailed.HttpCode(), model.ErrInternalDatabaseQueryFailedCode), false
	}
	return model.NoErr, result.RowsAffected == 1
}
