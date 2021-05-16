package typedErrors

import (
	"net/http"
)

var ForbiddenError = NewCustomError("Request - Forbidden operation", http.StatusForbidden)

func NewForbiddenErrorf(format string, args ...interface{}) error {
	return Wrapf(ForbiddenError, format, args...)
}
