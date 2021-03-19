package model

import (
	"cine-circle/internal/logger"
	"errors"
	"net/http"
)

const (
	ErrExternalSendRequestCode = "ERR_EXTERNAL_SEND_REQUEST"
	ErrExternalReceiveResponseCode = "ERR_EXTERNAL_RECEIVE_RESPONSE"
	ErrExternalReadBodyCode = "ERR_EXTERNAL_READ_BODY"

	ErrInternalDatabaseNilCode = "ERR_INTERNAL_DATABASE_NIL"
	ErrInternalDatabaseConnectionCode = "ERR_INTERNAL_DATABASE_CONNECTION"
	ErrInternalDatabaseCreationFailedCode = "ERR_INTERNAL_DATABASE_CREATION_FAILED"
	ErrInternalDatabaseQueryFailedCode = "ERR_INTERNAL_DATABASE_QUERY_FAILED"

	ErrInternalApiBadRequestCode = "ERR_INTERNAL_API_BAD_REQUEST"
	ErrInternalApiUnprocessableEntityCode = "ERR_INTERNAL_API_UNPROCESSABLE_ENTITY"
	ErrInternalApiNotFoundCode = "ERR_INTERNAL_API_NOT_FOUND"
)

var (
	ErrInternalDatabaseIsNil = NewCustomError(
		errors.New("got nil database when trying to connect"),
		http.StatusInternalServerError,
		ErrInternalDatabaseNilCode)

	ErrInternalApiBadRequest = NewCustomError(
		errors.New("request cannot be proceeded"),
		http.StatusBadRequest,
		ErrInternalApiBadRequestCode)

	ErrInternalApiUnprocessableEntity = NewCustomError(
		errors.New("cannot process entity"),
		http.StatusUnprocessableEntity,
		ErrInternalApiUnprocessableEntityCode)

	ErrInternalApiNotFound = NewCustomError(
		errors.New("entity cannot be found"),
		http.StatusNotFound,
		ErrInternalApiNotFoundCode)
)

type codeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CustomError struct {
	err error
	httpCode int
	code string
}

func (ce CustomError) Error() string {
	if ce.err != nil {
		return ce.err.Error()
	}
	return ""
}

func (ce CustomError) CodeError() codeError {
	return codeError{
		Code:    ce.code,
		Message: ce.Error(),
	}
}

func (ce CustomError) HttpCode() int {
	return ce.httpCode
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

func NewCustomError(err error, httpCode int, code string) CustomError {
	return CustomError{
		err:      err,
		httpCode: httpCode,
		code: code,
	}
}