package typedErrors

import (
	"net/http"
)

var ConfigurationError = NewCustomError("Application configuration", http.StatusInternalServerError)

func NewConfigurationErrorf(format string, args ...interface{}) error {
	return Wrapf(ConfigurationError, format, args...)
}
