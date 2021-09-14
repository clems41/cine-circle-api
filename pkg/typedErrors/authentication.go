package typedErrors

import (
	"net/http"
)

var AuthenticationError = NewCustomError("JWT - Authentication error", http.StatusUnauthorized)

func NewAuthenticationErrorf(format string, args ...interface{}) error {
	return Wrapf(AuthenticationError, format, args...)
}

func NewAuthenticationErrorWithCode(code string) error {
	return WrapCode(AuthenticationError, code)
}
