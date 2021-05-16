package userDom

import "cine-circle/internal/typedErrors"

const (
	errValidPasswordCode    = "PASSWORD_EMPTY"
	errValidEmailCode       = "EMAIL_EMPTY"
	errValidUsernameCode    = "USERNAME_EMPTY"
	errValidDisplayNameCode = "DISPLAY_NAME_EMPTY"
	errValidIDCode          = "ID_NULL"
	errValidOldPasswordCode = "OLD_PASSWORD_EMPTY"
	errValidNewPasswordCode = "NEW_PASSWORD_EMPTY"
	errValidGetCode         = "ALL_FIELDS_EMPTY"
	errValidKeywordCode     = "KEYWORD_IS_TOO_SHORT"
	errBadLoginPasswordCode = "BAD_LOGIN_PASSWORD"
)

var (
	errValidPassword    = typedErrors.NewBadRequestWithCode(errValidPasswordCode)
	errValidEmail       = typedErrors.NewBadRequestWithCode(errValidEmailCode)
	errValidUsername    = typedErrors.NewBadRequestWithCode(errValidUsernameCode)
	errValidDisplayName = typedErrors.NewBadRequestWithCode(errValidDisplayNameCode)
	errValidID          = typedErrors.NewBadRequestWithCode(errValidIDCode)
	errValidOldPassword = typedErrors.NewBadRequestWithCode(errValidOldPasswordCode)
	errValidNewPassword = typedErrors.NewBadRequestWithCode(errValidNewPasswordCode)
	errValidGet         = typedErrors.NewBadRequestWithCode(errValidGetCode)
	errValidKeyword     = typedErrors.NewBadRequestWithCode(errValidKeywordCode)
	errBadLoginPassword = typedErrors.NewAuthenticationErrorf(errBadLoginPasswordCode)
)
