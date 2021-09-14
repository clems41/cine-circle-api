package typedErrors

import (
	"net/http"
)

var NotImplementedError = NewCustomError("Server - This operation is not implemented", http.StatusNotImplemented)

func NewNotImplementedErrorf(format string, args ...interface{}) error {
	return Wrapf(NotImplementedError, format, args...)
}
