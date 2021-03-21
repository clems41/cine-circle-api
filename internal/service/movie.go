package service

import (
	"cine-circle/internal/model"
	"sort"
	"strings"
	"time"
)

func SortMovies(movies []model.Movie, sortParam string) []model.Movie {
	res := strings.Split(sortParam, ":")
	if len(res) != 2 {
		return nil
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
	return movies
}