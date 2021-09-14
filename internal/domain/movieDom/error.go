package movieDom

import (
	typedErrors2 "cine-circle/pkg/typedErrors"
)

const (
	errMovieNotFoundCode = "MOVIE_NOT_FOUND"
	errEmptyQueryCode = "EMPTY_QUERY"
)

var (
	errMovieNotFound = typedErrors2.NewNotFoundWithCode(errMovieNotFoundCode)
	errEmptyQuery = typedErrors2.NewBadRequestWithCode(errEmptyQueryCode)
)
