package recommendationDom

import (
	"cine-circle-api/pkg/customError"
)

const (
	errCircleNotFoundCode = "ERR_CIRCLE_NOT_FOUND"
	errUserNotFoundCode   = "ERR_USER_NOT_FOUND"
	errMovieNotFoundCode  = "ERR_MOVIE_NOT_FOUND"
)

var (
	errCircleNotFound = customError.NewNotFound().WrapCode(errCircleNotFoundCode)
	errUserNotFound   = customError.NewNotFound().WrapCode(errUserNotFoundCode)
	errMovieNotFound  = customError.NewNotFound().WrapCode(errMovieNotFoundCode)
)
