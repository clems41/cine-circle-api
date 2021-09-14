package typedErrors

import "fmt"

type TypedError interface {
	Error() string
	Type() string
	HttpStatus() int
	BusinessCode() string
}

type CustomError struct {
	typeString   string
	httpStatus   int
	wrapped      error
	businessCode string
}

func NewCustomError(typeString string, httpStatus int) CustomError {
	return CustomError{typeString: typeString, httpStatus: httpStatus}
}

func Wrapf(err CustomError, format string, args ...interface{}) error {
	return CustomError{
		typeString: fmt.Sprintf(format, args...),
		httpStatus: err.httpStatus,
		wrapped: err,
		businessCode: "",
	}
}

func WrapCode(err CustomError, code string) error {
	return CustomError{
		typeString: "",
		httpStatus: err.httpStatus,
		wrapped: err,
		businessCode: code,
	}
}

func (err CustomError) Error() string {
	return err.Type()
}

func (err CustomError) Type() string {
	return err.typeString
}

func (err CustomError) HttpStatus() int {
	return err.httpStatus
}

func (err CustomError) Unwrap() error {
	return err.wrapped
}

func (err CustomError) BusinessCode() string {
	return err.businessCode
}
