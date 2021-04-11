package service

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/model"
	"cine-circle/internal/repository"
	"cine-circle/internal/typedErrors"
	"sort"
)

func AddRating(rating model.Rating, movieId, username string) (typedErrors.CustomError, model.Rating) {
	if !omdb.MovieExists(movieId) {
		return typedErrors.ErrRepositoryResourceNotFound, rating
	}
	if rating.Value > model.RatingBoundMax || rating.Value < model.RatingBoundMin {
		return typedErrors.ErrApiUnprocessableEntity, rating
	}
	err, user := GetUser("username = ?", username)
	if err.IsNotNil() {
		return err, rating
	}
	rating.UserID = user.ID
	rating.MovieID = movieId
	rating.Source = model.CineCircleSource
	rating.Username = user.Username
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err, rating
	}
	defer db.Close()
	err = db.CreateOrUpdate(&model.Rating{}, &rating, "movie_id = ? AND user_id = ?", rating.MovieID, rating.UserID)
	return err, rating
}

func AddUserRatings(username string, movie *model.Movie) typedErrors.CustomError {
	if username == "" {
		return typedErrors.NoErr
	}
	var ratings []model.Rating
	err, user := GetUser("username = ?", username)
	var circlesId []uint
	var usersId []uint
	db, err := repository.OpenConnection()
	if err.IsNotNil() {
		return err
	}
	defer db.Close()
	result := db.DB().Table("user_circle").Select("circle_id").Find(&circlesId, "user_id = ?", user.ID)
	if result.RowsAffected > 0 {
		db.DB().Table("user_circle").Select("user_id").Find(&usersId, "circle_id IN (?)", circlesId)
	} else {
		usersId = append(usersId, user.ID)
	}
	result = db.DB().Find(&ratings, "user_id IN (?) AND movie_id = ?", filterUsersId(usersId), movie.ID)
	if result.RowsAffected > 0 {
		sortRatings(ratings)
		movie.UserRatings = append(movie.UserRatings, ratings...)
	}
	return err
}

func filterUsersId(ids []uint) []uint {
	resultMap := make(map[uint]uint)
	var result []uint
	for _, id := range ids {
		resultMap[id] = id
	}
	for newId, _ := range resultMap {
		result = append(result, newId)
	}
	return result
}

func sortRatings(ratings []model.Rating) {
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].UpdatedAt.After(ratings[j].UpdatedAt)
	})
}