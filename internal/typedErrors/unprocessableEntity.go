package typedErrors

import "net/http"

var UnprocessableEntityError = NewCustomError("Entity is not processable", http.StatusUnprocessableEntity)

func NewUnprocessableEntityErrorf(format string, args ...interface{}) error {
	return Wrapf(UnprocessableEntityError, format, args...)
}

func NewUnprocessableEntityWithCode(code string) error {
	return WrapCode(UnprocessableEntityError, code)
}
