package typedErrors

import "net/http"

var ConflictError = NewCustomError("Request - Conflict", http.StatusConflict)

func NewConflictErrorf(format string, args ...interface{}) error {
	return Wrapf(BadRequestError, format, args...)
}

func NewConflictWithCode(code string) error {
	return WrapCode(ConflictError, code)
}
