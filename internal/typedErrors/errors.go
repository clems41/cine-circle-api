package typedErrors

import (
	"cine-circle/internal/logger"
	"errors"
	"net/http"
)

const (
	errExternalSendRequestCode = "ERR_EXTERNAL_SEND_REQUEST_CODE"
	errExternalReceiveResponseCode = "ERR_EXTERNAL_RECEIVE_RESPONSE_CODE"
	errExternalReadBodyCode = "ERR_EXTERNAL_READ_BODY_CODE"

	errRepositoryNilCode = "ERR_INTERNAL_DATABASE_NIL_CODE"
	errRepositoryResourceNotFoundCode = "ERR_INTERNAL_DATABASE_RESOURCE_NOT_FOUND_CODE"
	errRepositoryConnectionFailedCode = "ERR_INTERNAL_DATABASE_CONNECTION_FAILED_CODE"
	errRepositoryCreationFailedCode = "ERR_INTERNAL_DATABASE_CREATION_FAILED_CODE"
	errRepositoryQueryFailedCode = "ERR_INTERNAL_DATABASE_QUERY_FAILED_CODE"

	errApiBadRequestCode = "ERR_INTERNAL_API_BAD_REQUEST_CODE"
	errApiUnprocessableEntityCode = "ERR_INTERNAL_API_UNPROCESSABLE_ENTITY_CODE"
	errApiUserBadCredentialsCode = "ERR_INTERNAL_API_USER_BAD_CREDENTIALS_CODE"
	errApiUserCredentialsNotFoundCode = "ERR_INTERNAL_API_USER_CREDENTIALS_NOT_FOUND_CODE"

	errServiceMissingMandatoryFieldsCode = "ERR_INTERNAL_SERVICE_MISSING_MANDATORY_FIELDS_CODE"
	errServiceBadFormatMandatoryFieldsCode = "ERR_INTERNAL_SERVICE_BAD_FORMAT_MANDATORY_FIELDS_CODE"
	errServiceFieldShouldBeUniqueCode = "ERR_INTERNAL_SERVICE_FIELD_SHOULD_BE_UNIQUE_CODE"
)

var (
	NoErr = NewCustomError(nil,	http.StatusOK,"")

	ErrRepositoryIsNil = NewCustomError(
		errors.New("got nil repository when trying to connect"),
		http.StatusInternalServerError,
		errRepositoryNilCode)
	ErrRepositoryResourceNotFound = NewCustomError(
		errors.New("resource cannot be found in repository"),
		http.StatusNotFound,
		errRepositoryResourceNotFoundCode)
	ErrRepositoryConnectionFailed = NewCustomError(
		errors.New("connection to repository failed"),
		http.StatusInternalServerError,
		errRepositoryConnectionFailedCode)
	ErrRepositoryQueryFailed = NewCustomError(
		errors.New("query sent to repository failed"),
		http.StatusInternalServerError,
		errRepositoryQueryFailedCode)

	ErrApiBadRequest = NewCustomError(
		errors.New("request cannot be proceeded"),
		http.StatusBadRequest,
		errApiBadRequestCode)
	ErrApiUnprocessableEntity = NewCustomError(
		errors.New("cannot process entity"),
		http.StatusUnprocessableEntity,
		errApiUnprocessableEntityCode)
	ErrApiUserBadCredentials = NewCustomError(
		errors.New("cannot authenticate user with these credentials"),
		http.StatusUnauthorized,
		errApiUserBadCredentialsCode)
	ErrApiUserCredentialsNotFound = NewCustomError(
		errors.New("user credentials not found"),
		http.StatusUnauthorized,
		errApiUserCredentialsNotFoundCode)

	ErrServiceMissingMandatoryFields = NewCustomError(
		errors.New("missing mandatory fields for creating resource"),
		http.StatusUnprocessableEntity,
		errServiceMissingMandatoryFieldsCode)
	ErrServiceBadFormatMandatoryFields = NewCustomError(
		errors.New("bad format for at least one mandatory fields"),
		http.StatusUnprocessableEntity,
		errServiceBadFormatMandatoryFieldsCode)
	ErrServiceFieldShouldBeUnique = NewCustomError(
		errors.New("field should be unique, already used for other resource"),
		http.StatusUnprocessableEntity,
		errServiceFieldShouldBeUniqueCode)
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

// Repository errors

func NewRepositoryQueryFailedError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errRepositoryQueryFailedCode)
}
func NewRepositoryConnectionFailedError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errRepositoryConnectionFailedCode)
}
func NewRepositoryCreationFailedError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errRepositoryCreationFailedCode)
}
func NewRepositoryNilError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errRepositoryNilCode)
}
func NewRepositoryResourceNotFoundError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errRepositoryResourceNotFoundCode)
}

// API errors

func NewApiBadCredentialsError(err error) CustomError {
	return NewCustomError(err, http.StatusBadRequest, errApiUserBadCredentialsCode)
}
func NewApiCredentialsNotFoundError(err error) CustomError {
	return NewCustomError(err, http.StatusBadRequest, errApiUserCredentialsNotFoundCode)
}
func NewApiBadRequestError(err error) CustomError {
	return NewCustomError(err, http.StatusBadRequest, errApiBadRequestCode)
}
func NewApiUnprocessableEntityError(err error) CustomError {
	return NewCustomError(err, http.StatusUnprocessableEntity, errApiUnprocessableEntityCode)
}

// Service errors

func NewServiceMissingMandatoryFieldsError(err error) CustomError {
	return NewCustomError(err, http.StatusBadRequest, errServiceMissingMandatoryFieldsCode)
}
func NewServiceBadFormatMandatoryFieldsError(err error) CustomError {
	return NewCustomError(err, http.StatusUnprocessableEntity, errServiceBadFormatMandatoryFieldsCode)
}
func NewServiceFieldsShouldBeUniqueError(err error) CustomError {
	return NewCustomError(err, http.StatusUnprocessableEntity, errServiceFieldShouldBeUniqueCode)
}

// External errors

func NewExternalSendRequestError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errExternalSendRequestCode)
}
func NewExternalReadBodyError(err error) CustomError {
	return NewCustomError(err, http.StatusInternalServerError, errExternalReadBodyCode)
}