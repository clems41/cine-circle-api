package typedErrors

import (
	"net/http"
)

var GoneError = NewCustomError("Request - Entity has been deleted", http.StatusGone)

func NewGoneErrorf(format string, args ...interface{}) error {
	return Wrapf(GoneError, format, args...)
}
