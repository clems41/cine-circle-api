package service

import (
	"cine-circle/internal/model"
	"sort"
	"strings"
	"time"
)

func SortMovies(movies []model.Movie, sortParam string) model.CustomError {
	res := strings.Split(sortParam, ":")
	if len(res) != 2 {
		return model.ErrInternalApiBadRequest
	}
	field := res[0]
	asc := res[1] == "asc"
	switch field {
	case "title":
		sort.SliceStable(movies, func(i, j int) bool {
			return movies[i].Title < movies[j].Title == asc
		})
	default:
		sort.SliceStable(movies, func(i, j int) bool {
			var firstPostedDate, secondPostedDate time.Time
			for _, rating := range movies[i].UserRatings {
				if asc == rating.UpdatedAt.Before(firstPostedDate) || firstPostedDate.IsZero() {
					firstPostedDate = rating.UpdatedAt
				}
			}
			for _, rating := range movies[j].UserRatings {
				if asc == rating.UpdatedAt.Before(firstPostedDate) || secondPostedDate.IsZero() {
					secondPostedDate = rating.UpdatedAt
				}
			}
			return firstPostedDate.Before(secondPostedDate) == asc
		})
	}
	return model.NoErr
}

func MergeMovies(movies []model.Movie) (model.CustomError, []model.Movie) {
	moviesMerged := make(map[string]model.Movie)
	moviesRatings := make(map[string][]model.Rating)
	var result []model.Movie
	for _, movie := range movies {
		moviesMerged[movie.ID] = movie
		moviesRatings[movie.ID] = append(moviesRatings[movie.ID], movie.UserRatings...)
	}
	for movieId, movieMerged := range moviesMerged {
		movieMerged.UserRatings = moviesRatings[movieId]
		result = append(result, movieMerged)
	}
	return model.NoErr, result
}