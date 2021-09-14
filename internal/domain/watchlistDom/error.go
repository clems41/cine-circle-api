package watchlistDom

import (
	"cine-circle/pkg/typedErrors"
)

const (
	errMovieNotFoundCode    = "MOVIE_NOT_FOUND"
)

var (
	errMovieNotFound    = typedErrors.NewNotFoundWithCode(errMovieNotFoundCode)
)
