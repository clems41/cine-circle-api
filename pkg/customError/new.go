package customError

import "net/http"

func NewBadRequest() *CustomError {
	return &CustomError{
		httpStatus: http.StatusBadRequest,
	}
}

func NewUnauthorized() *CustomError {
	return &CustomError{
		httpStatus: http.StatusUnauthorized,
	}
}

func NewForbidden() *CustomError {
	return &CustomError{
		httpStatus: http.StatusForbidden,
	}
}

func NewUnprocessableEntity() *CustomError {
	return &CustomError{
		httpStatus: http.StatusUnprocessableEntity,
	}
}

func NewInternalServer() *CustomError {
	return &CustomError{
		httpStatus: http.StatusInternalServerError,
	}
}

func NewNotFound() *CustomError {
	return &CustomError{
		httpStatus: http.StatusNotFound,
	}
}
