package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/database"
	"cine-circle/internal/model"
	"net/http"
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
	db.DB().Exec("INSERT INTO watchlist(user_id, movie_id) VALUES (?, ?)", user.ID, movieId)
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
	db.DB().Exec("DELETE FROM watchlist WHERE user_id = ? AND movie_id = ?", user.ID, movieId)
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
	var moviesId []string
	db.DB().Table("watchlist").Select("movie_id").Find(&moviesId, "user_id = ?", user.ID)
	var result model.MovieSearch
	for _, movieId := range moviesId {
		err, movie := omdb.FindMovieByID(movieId)
		if err.IsNotNil() {
			return err, model.MovieSearch{}
		}
		result.Search = append(result.Search, movie.MovieShort())
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
	var moviesId []string
	result := db.DB().Table("watchlist").Select("movie_id").Find(&moviesId, "user_id = ? AND movie_id = ?", user.ID, movieId)
	if result.Error != nil {
		return model.NewCustomError(result.Error, http.StatusBadRequest, model.ErrInternalDatabaseQueryFailedCode), false
	}
	return model.NoErr, result.RowsAffected == 1
}
