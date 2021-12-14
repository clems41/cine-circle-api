package userDom

import "cine-circle-api/pkg/utils/httpUtils"

/* Routes */

const (
	basePath = "/v1/login"

	signInPath                 = "/sign-in"
	signUpPath                 = "/sign-up"
	sendConfirmationEmailPath  = "/send-confirmation-email"
	confirmEmailPath           = "/confirm-email"
	sendResetPasswordEmailPath = "/send-reset-password-email"
	resetPasswordPath          = "/reset-password"

	ownInfoPath        = "/me"
	updatePasswordPath = "/password"
)

/* Path parameters */

const (
	loginPathParameter httpUtils.PathParameter = "login"
)
