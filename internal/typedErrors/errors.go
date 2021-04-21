package typedErrors

import (
	"cine-circle/internal/logger"
	"errors"
	"fmt"
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
	errServiceGeneralErrorCode = "ERR_INTERNAL_SERVICE_ERROR_CODE"
)

var (
	ErrRepositoryIsNil = NewCustomError(
		http.StatusInternalServerError,
		errRepositoryNilCode,
		errors.New("got nil repository when trying to connect"))
	ErrRepositoryResourceNotFound = NewCustomError(
		http.StatusNotFound,
		errRepositoryResourceNotFoundCode,
		errors.New("resource cannot be found in repository"))
	ErrRepositoryConnectionFailed = NewCustomError(
		http.StatusInternalServerError,
		errRepositoryConnectionFailedCode,
	errors.New("connection to repository failed"))
	ErrRepositoryQueryFailed = NewCustomError(
		http.StatusInternalServerError,
		errRepositoryQueryFailedCode,
		errors.New("query sent to repository failed"))

	ErrApiBadRequest = NewCustomError(
		http.StatusBadRequest,
		errApiBadRequestCode,
		errors.New("request cannot be proceeded"))
	ErrApiUnprocessableEntity = NewCustomError(
		http.StatusUnprocessableEntity,
		errApiUnprocessableEntityCode,
		errors.New("cannot process entity"))
	ErrApiUserBadCredentials = NewCustomError(
		http.StatusUnauthorized,
		errApiUserBadCredentialsCode,
		errors.New("cannot authenticate user with these credentials"))
	ErrApiUserCredentialsNotFound = NewCustomError(
		http.StatusUnauthorized,
		errApiUserCredentialsNotFoundCode,
		errors.New("user credentials not found"))

	ErrServiceMissingMandatoryFields = NewCustomError(
		http.StatusUnprocessableEntity,
		errServiceMissingMandatoryFieldsCode,
		errors.New("missing mandatory fields for creating resource"))
	ErrServiceBadFormatMandatoryFields = NewCustomError(
		http.StatusUnprocessableEntity,
		errServiceBadFormatMandatoryFieldsCode,
		errors.New("bad format for at least one mandatory fields"))
	ErrServiceFieldShouldBeUnique = NewCustomError(
		http.StatusUnprocessableEntity,
		errServiceFieldShouldBeUniqueCode,
		errors.New("field should be unique, already used for other resource"))
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

func NewCustomErrorf(httpCode int, code string, format string, args ...interface{}) CustomError {
	newErr := errors.New(fmt.Sprintf(format, args...))
	return CustomError{
		err:      fmt.Errorf("decompress %s: %w", newErr.Error(), newErr),
		httpCode: httpCode,
		code: code,
	}
}

func NewCustomError(httpCode int, code string, err error) CustomError {
	return CustomError{
		err:      err,
		httpCode: httpCode,
		code: code,
	}
}

// Repository errors

func NewRepositoryQueryFailedErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusNotFound, errRepositoryQueryFailedCode, format, args...)
}
func NewRepositoryConnectionFailedErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errRepositoryConnectionFailedCode, format, args...)
}
func NewRepositoryCreationFailedErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errRepositoryCreationFailedCode, format, args...)
}
func NewRepositoryNilErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errRepositoryNilCode, format, args...)
}
func NewRepositoryResourceNotFoundErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errRepositoryResourceNotFoundCode, format, args...)
}

// API errors

func NewApiBadCredentialsErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusBadRequest, errApiUserBadCredentialsCode, format, args...)
}
func NewApiCredentialsNotFoundErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusBadRequest, errApiUserCredentialsNotFoundCode, format, args...)
}
func NewApiBadRequestErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusBadRequest, errApiBadRequestCode, format, args...)
}
func NewApiUnprocessableEntityErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusUnprocessableEntity, errApiUnprocessableEntityCode, format, args...)
}

// Service errors

func NewServiceMissingMandatoryFieldsErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusBadRequest, errServiceMissingMandatoryFieldsCode, format, args...)
}
func NewServiceBadFormatMandatoryFieldsErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusUnprocessableEntity, errServiceBadFormatMandatoryFieldsCode, format, args...)
}
func NewServiceFieldsShouldBeUniqueErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusUnprocessableEntity, errServiceFieldShouldBeUniqueCode, format, args...)
}
func NewServiceGeneralErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errServiceGeneralErrorCode, format, args...)
}

// External errors

func NewExternalSendRequestErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errExternalSendRequestCode, format, args...)
}
func NewExternalReadBodyErrorf(format string, args ...interface{}) CustomError {
	return NewCustomErrorf(http.StatusInternalServerError, errExternalReadBodyCode, format, args...)
}