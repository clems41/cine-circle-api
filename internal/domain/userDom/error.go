package userDom

import (
	typedErrors2 "cine-circle/pkg/typedErrors"
)

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
	errValidPassword    = typedErrors2.NewBadRequestWithCode(errValidPasswordCode)
	errValidEmail       = typedErrors2.NewBadRequestWithCode(errValidEmailCode)
	errValidUsername    = typedErrors2.NewBadRequestWithCode(errValidUsernameCode)
	errValidDisplayName = typedErrors2.NewBadRequestWithCode(errValidDisplayNameCode)
	errValidID          = typedErrors2.NewBadRequestWithCode(errValidIDCode)
	errValidOldPassword = typedErrors2.NewBadRequestWithCode(errValidOldPasswordCode)
	errValidNewPassword = typedErrors2.NewBadRequestWithCode(errValidNewPasswordCode)
	errValidGet         = typedErrors2.NewBadRequestWithCode(errValidGetCode)
	errValidKeyword     = typedErrors2.NewBadRequestWithCode(errValidKeywordCode)
	errBadLoginPassword = typedErrors2.NewAuthenticationErrorf(errBadLoginPasswordCode)
)
