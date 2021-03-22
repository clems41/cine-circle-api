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
			for _, rating := range movies[i].Ratings {
				if rating.Source == model.RatingSourceCineCircle {
					firstPostedDate = rating.PostedDate
				}
			}
			for _, rating := range movies[j].Ratings {
				if rating.Source == model.RatingSourceCineCircle {
					secondPostedDate = rating.PostedDate
				}
			}
			return firstPostedDate.Before(secondPostedDate) == asc
		})
	}
	return model.NoErr
}

func MergeMovies(movies []model.Movie, usersId []uint) (model.CustomError, []model.Movie) {
	moviesMerged := make(map[string]model.Movie)
	moviesRatings := make(map[string][]model.MovieRating)
	var result []model.Movie
	for _, movie := range movies {
		if _, exists := moviesMerged[movie.Imdbid]; !exists {
			moviesMerged[movie.Imdbid] = movie
			moviesRatings[movie.Imdbid] = append(moviesRatings[movie.Imdbid], movie.Ratings...)
		} else {
			userRatingIdx := -1
			for ratingIdx, rating := range movie.Ratings {
				if rating.Source == model.RatingSourceCineCircle {
					userRatingIdx = ratingIdx
				}
			}
			if userRatingIdx >= 0 {
				moviesRatings[movie.Imdbid] = append(moviesRatings[movie.Imdbid], movie.Ratings[userRatingIdx])
			}
		}
	}
	for movieId, movieMerged := range moviesMerged {
		movieMerged.Ratings = moviesRatings[movieId]
		result = append(result, movieMerged)
	}
	return model.NoErr, result
}