package model

import (
	"cine-circle/internal/logger"
	"errors"
	"net/http"
)

const (
	ErrExternalSendRequestCode = "ERR_EXTERNAL_SEND_REQUEST_CODE"
	ErrExternalReceiveResponseCode = "ERR_EXTERNAL_RECEIVE_RESPONSE_CODE"
	ErrExternalReadBodyCode = "ERR_EXTERNAL_READ_BODY_CODE"

	ErrInternalDatabaseNilCode = "ERR_INTERNAL_DATABASE_NIL_CODE"
	ErrInternalDatabaseResourceNotFoundCode = "ERR_INTERNAL_DATABASE_RESOURCE_NOT_FOUND_CODE"
	ErrInternalDatabaseConnectionFailedCode = "ERR_INTERNAL_DATABASE_CONNECTION_FAILED_CODE"
	ErrInternalDatabaseCreationFailedCode = "ERR_INTERNAL_DATABASE_CREATION_FAILED_CODE"
	ErrInternalDatabaseQueryFailedCode = "ERR_INTERNAL_DATABASE_QUERY_FAILED_CODE"

	ErrInternalApiBadRequestCode = "ERR_INTERNAL_API_BAD_REQUEST_CODE"
	ErrInternalApiUnprocessableEntityCode = "ERR_INTERNAL_API_UNPROCESSABLE_ENTITY_CODE"
	ErrInternalApiUserBadCredentialsCode = "ERR_INTERNAL_API_USER_BAD_CREDENTIALS_CODE"
	ErrInternalApiUserCredentialsNotFoundCode = "ERR_INTERNAL_API_USER_CREDENTIALS_NOT_FOUND_CODE"

	ErrInternalServiceMissingMandatoryFieldsCode = "ERR_INTERNAL_SERVICE_MISSING_MANDATORY_FIELDS_CODE"
	ErrInternalServiceBadFormatMandatoryFieldsCode = "ERR_INTERNAL_SERVICE_BAD_FORMAT_MANDATORY_FIELDS_CODE"
	ErrInternalServiceFieldShouldBeUniqueCode = "ERR_INTERNAL_SERVICE_FIELD_SHOULD_BE_UNIQUE_CODE"
)

var (
	NoErr = NewCustomError(nil,	http.StatusOK,"")

	ErrInternalDatabaseIsNil = NewCustomError(
		errors.New("got nil database when trying to connect"),
		http.StatusInternalServerError,
		ErrInternalDatabaseNilCode)
	ErrInternalDatabaseResourceNotFound = NewCustomError(
		errors.New("resource cannot be found in database"),
		http.StatusNotFound,
		ErrInternalDatabaseResourceNotFoundCode)
	ErrInternalDatabaseConnectionFailed = NewCustomError(
		errors.New("connection to database failed"),
		http.StatusInternalServerError,
		ErrInternalDatabaseConnectionFailedCode)
	ErrInternalDatabaseQueryFailed = NewCustomError(
		errors.New("query sent to database failed"),
		http.StatusInternalServerError,
		ErrInternalDatabaseQueryFailedCode)

	ErrInternalApiBadRequest = NewCustomError(
		errors.New("request cannot be proceeded"),
		http.StatusBadRequest,
		ErrInternalApiBadRequestCode)
	ErrInternalApiUnprocessableEntity = NewCustomError(
		errors.New("cannot process entity"),
		http.StatusUnprocessableEntity,
		ErrInternalApiUnprocessableEntityCode)
	ErrInternalApiUserBadCredentials = NewCustomError(
		errors.New("cannot authenticate user with these credentials"),
		http.StatusUnauthorized,
		ErrInternalApiUserBadCredentialsCode)
	ErrInternalApiUserCredentialsNotFound = NewCustomError(
		errors.New("user credentials not found"),
		http.StatusUnauthorized,
		ErrInternalApiUserCredentialsNotFoundCode)

	ErrInternalServiceMissingMandatoryFields = NewCustomError(
		errors.New("missing mandatory fields for creating resource"),
		http.StatusUnprocessableEntity,
		ErrInternalServiceMissingMandatoryFieldsCode)
	ErrInternalServiceBadFormatMandatoryFields = NewCustomError(
		errors.New("bad format for at least one mandatory fields"),
		http.StatusUnprocessableEntity,
		ErrInternalServiceBadFormatMandatoryFieldsCode)
	ErrInternalServiceFieldShouldBeUnique = NewCustomError(
		errors.New("field should be unique, already used for other resource"),
		http.StatusUnprocessableEntity,
		ErrInternalServiceFieldShouldBeUniqueCode)
)

type codeError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
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