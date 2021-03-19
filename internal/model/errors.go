package model

import (
	"cine-circle/internal/logger"
	"errors"
	"net/http"
)

var (
	ErrInternalDatabaseIsNil = NewCustomError(
		errors.New("got nil database when trying to connect"),
		http.StatusInternalServerError)
)

type CustomError struct {
	err error
	httpCode int
}

func (ce CustomError) Error() string {
	if ce.err != nil {
		return ce.err.Error()
	}
	return ""
}

func (ce CustomError) Print() {
	logger.Sugar.Errorf("Erorr occurs : %s", ce.Error())
}

func (ce CustomError) IsNotNil() bool {
	isNotNil := ce.err != nil
	if isNotNil {
		ce.Print()
	}
	return isNotNil
}

func NewCustomError(err error, httpCode int) CustomError {
	return CustomError{
		err:      err,
		httpCode: httpCode,
	}
}