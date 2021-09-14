package circleDom

import (
	typedErrors2 "cine-circle/pkg/typedErrors"
)

const (
	errNameEmptyCode        = "NAME_EMPTY"
	errDescriptionEmptyCode = "DESCRIPTION_EMPTY"
	errNoFieldsProvidedCode = "NO_FIELDS_PROVIDED"
	errIdNullCode           = "ID_NULL"
	errNotAuthorizedCode    = "USER_NOT_AUTHORIZED"
	errUserNotFoundCode     = "USER_NOT_FOUND"
)

var (
	errNameEmpty        = typedErrors2.NewBadRequestWithCode(errNameEmptyCode)
	errDescriptionEmpty = typedErrors2.NewBadRequestWithCode(errDescriptionEmptyCode)
	errNoFieldsProvided = typedErrors2.NewBadRequestWithCode(errNoFieldsProvidedCode)
	errIdNull           = typedErrors2.NewBadRequestWithCode(errIdNullCode)
	errNotAuthorized    = typedErrors2.NewAuthenticationErrorWithCode(errNotAuthorizedCode)
	errUserNotFound     = typedErrors2.NewNotFoundWithCode(errUserNotFoundCode)
)
