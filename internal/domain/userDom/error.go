package userDom

import "cine-circle-api/pkg/customError"

const (
	errBadCredentialsCode      = "ERR_BAD_CREDENTIALS"
	errUserUnauthorizedCode    = "ERR_USER_UNAUTHORIZED"
	errUserCannotBeDeletedCode = "ERR_USER_CANNOT_BE_DELETED"
)

var (
	errInvalidAuthorizationHeader = customError.NewUnauthorized().WrapErrorf("cannot retrieved login and password from Authorization header")
	errBadCredentials             = customError.NewUnauthorized().WrapCode(errBadCredentialsCode)
	errUserUnauthorized           = customError.NewUnauthorized().WrapCode(errUserUnauthorizedCode)
	errUserCannotBeDeleted        = customError.NewInternalServer().WrapCode(errUserCannotBeDeletedCode)
)
