package typedErrors

import (
	"net/http"
)

var NotFoundError = NewCustomError("Request - Entity not found", http.StatusNotFound)

func NewNotFoundErrorf(format string, args ...interface{}) error {
	return Wrapf(NotFoundError, format, args...)
}
