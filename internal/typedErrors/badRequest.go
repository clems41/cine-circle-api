package typedErrors

import (
	"net/http"
)

var BadRequestError = NewCustomError("Request - Bad service use", http.StatusBadRequest)

func NewBadRequestErrorf(format string, args ...interface{}) error {
	return Wrapf(BadRequestError, format, args...)
}

func NewBadRequestWithCode(code string) error {
	return WrapCode(BadRequestError, code)
}
