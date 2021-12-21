package userDom

import "cine-circle-api/pkg/customError"

const (
	errUsernameAlreadyExistsCode = "ERR_USERNAME_ALREADY_EXISTS"
	errEmailAlreadyExistsCode    = "ERR_EMAIL_ALREADY_EXISTS"
	errBadCredentialsCode        = "ERR_BAD_CREDENTIALS"
	errUserUnauthorizedCode      = "ERR_USER_UNAUTHORIZED"
	errUserCannotBeDeletedCode   = "ERR_USER_CANNOT_BE_DELETED"
)

var (
	errUsernameAlreadyExists      = customError.NewBadRequest().WrapCode(errUsernameAlreadyExistsCode)
	errEmailAlreadyExists         = customError.NewBadRequest().WrapCode(errEmailAlreadyExistsCode)
	errInvalidAuthorizationHeader = customError.NewUnauthorized().WrapErrorf("cannot retrieved login and password from Authorization header")
	errBadCredentials             = customError.NewUnauthorized().WrapCode(errBadCredentialsCode)
	errUserUnauthorized           = customError.NewUnauthorized().WrapCode(errUserUnauthorizedCode)
	errUserCannotBeDeleted        = customError.NewInternalServer().WrapCode(errUserCannotBeDeletedCode)
)
