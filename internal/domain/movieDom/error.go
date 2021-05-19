package movieDom

import "cine-circle/internal/typedErrors"

const (
	errMovieNotFoundCode = "MOVIE_NOT_FOUND"
)

var (
	errMovieNotFound = typedErrors.NewNotFoundWithCode(errMovieNotFoundCode)
)
