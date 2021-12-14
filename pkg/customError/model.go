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
	str := fmt.Sprintf("%s [%d]", ce.code, ce.httpStatus)
	if ce.error != nil {
		str += fmt.Sprintf(" - %s", ce.error.Error())
	}
	return str
}

func (ce *CustomError) HttpStatus() int {
	return ce.httpStatus
}

func (ce *CustomError) Code() string {
	return ce.code
}
