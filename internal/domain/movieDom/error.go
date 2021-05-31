package movieDom

import "cine-circle/internal/typedErrors"

const (
	errMovieNotFoundCode = "MOVIE_NOT_FOUND"
	errEmptyQueryCode = "EMPTY_QUERY"
)

var (
	errMovieNotFound = typedErrors.NewNotFoundWithCode(errMovieNotFoundCode)
	errEmptyQuery = typedErrors.NewBadRequestWithCode(errEmptyQueryCode)
)
