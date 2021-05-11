package userDom

import "cine-circle/internal/typedErrors"

var (
	errValidPassword    = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Password is empty")
	errValidEmail       = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Email is empty")
	errValidUsername    = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Username is empty")
	errValidDisplayName = typedErrors.NewServiceMissingMandatoryFieldsErrorf("DisplayName is empty")
	errValidID          = typedErrors.NewServiceMissingMandatoryFieldsErrorf("UserID is empty")
	errValidOldPassword = typedErrors.NewServiceMissingMandatoryFieldsErrorf("OldPassword is empty")
	errValidNewPassword = typedErrors.NewServiceMissingMandatoryFieldsErrorf("NewPassword is empty")
	errValidGet         = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Need at least one field to get user")
	errValidUpdate      = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Need at least one field to update user")
	errValidKeyword     = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Need at least 3 chars for searching for users")
	errBadLoginPassword = typedErrors.NewApiBadCredentialsErrorf("wrong login/password")
)
