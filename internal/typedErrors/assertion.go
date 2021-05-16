package typedErrors

import (
	"net/http"
)

var AssertionError = NewCustomError("Assertion error", http.StatusInternalServerError)

func NewAssertionErrorf(format string, args ...interface{}) error {
	return Wrapf(AssertionError, format, args...)
}
