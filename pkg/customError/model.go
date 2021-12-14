package customError

import "fmt"

type CustomError struct {
	error      error
	code       string
	httpStatus int
}

func (ce *CustomError) WrapCode(code string) *CustomError {
	ce.code = code
	return ce
}

func (ce *CustomError) WrapError(err error) *CustomError {
	ce.error = err
	return ce
}

func (ce *CustomError) WrapErrorf(format string, a ...interface{}) *CustomError {
	ce.error = fmt.Errorf(format, a...)
	return ce
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s [%d] - %s", ce.code, ce.httpStatus, ce.error.Error())
}

func (ce *CustomError) HttpStatus() int {
	return ce.httpStatus
}

func (ce *CustomError) Code() string {
	return ce.code
}
