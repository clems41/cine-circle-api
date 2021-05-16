package circleDom

import "cine-circle/internal/typedErrors"

const (
	errNameEmptyCode = "NAME_EMPTY"
	errDescriptionEmptyCode = "DESCRIPTION_EMPTY"
	errUsersEmptyCode = "USERS_EMPTY"
	errNoFieldsProvidedCode = "NO_FIELDS_PROVIDED"
	errIdNullCode = "ID_NULL"
	errNotAuthorizedCode = "USER_NOT_AUTHORIZED"
)

var (
	errNameEmpty = typedErrors.NewBadRequestWithCode(errNameEmptyCode)
	errDescriptionEmpty = typedErrors.NewBadRequestWithCode(errDescriptionEmptyCode)
	errUsersEmpty = typedErrors.NewBadRequestWithCode(errUsersEmptyCode)
	errNoFieldsProvided = typedErrors.NewBadRequestWithCode(errNoFieldsProvidedCode)
	errIdNull = typedErrors.NewBadRequestWithCode(errIdNullCode)
	errNotAuthorized = typedErrors.NewAuthenticationErrorWithCode(errNotAuthorizedCode)
)
