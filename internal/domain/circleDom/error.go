package circleDom

import "cine-circle/internal/typedErrors"

const (
	errNameEmptyCode        = "NAME_EMPTY"
	errDescriptionEmptyCode = "DESCRIPTION_EMPTY"
	errNoFieldsProvidedCode = "NO_FIELDS_PROVIDED"
	errIdNullCode           = "ID_NULL"
	errNotAuthorizedCode    = "USER_NOT_AUTHORIZED"
	errUserNotFoundCode     = "USER_NOT_FOUND"
)

var (
	errNameEmpty        = typedErrors.NewBadRequestWithCode(errNameEmptyCode)
	errDescriptionEmpty = typedErrors.NewBadRequestWithCode(errDescriptionEmptyCode)
	errNoFieldsProvided = typedErrors.NewBadRequestWithCode(errNoFieldsProvidedCode)
	errIdNull           = typedErrors.NewBadRequestWithCode(errIdNullCode)
	errNotAuthorized    = typedErrors.NewAuthenticationErrorWithCode(errNotAuthorizedCode)
	errUserNotFound     = typedErrors.NewNotFoundWithCode(errUserNotFoundCode)
)
