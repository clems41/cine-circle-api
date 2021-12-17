package userDom

import "cine-circle-api/pkg/utils/httpUtils"

/* Http parameters */

var (
	loginParameter = httpUtils.Parameter{
		Name:            "login",
		Description:     "Login (username or email) to use to reset password",
		DefaultValueStr: "",
		DataType:        "string",
		Required:        true,
	}
	keywordQueryParameter = httpUtils.Parameter{
		Name:            "keyword",
		Description:     "Keyword to search among users on firstname, lastname and username fields",
		DefaultValueStr: "",
		DataType:        "string",
		Required:        true,
	}
	pageQueryParameter = httpUtils.Parameter{
		Name:            "page",
		Description:     "Page number to return",
		DefaultValueStr: "1",
		DataType:        "int",
		Required:        false,
	}
	pageSizeQueryParameter = httpUtils.Parameter{
		Name:            "pageSize",
		Description:     "Number of element by page",
		DefaultValueStr: "10",
		DataType:        "int",
		Required:        false,
	}
	defaultSearchQueryParameterValues = map[string]string{
		pageQueryParameter.Name:     pageQueryParameter.DefaultValueStr,
		pageSizeQueryParameter.Name: pageSizeQueryParameter.DefaultValueStr,
	}
)

/* Routes */

const (
	basePath = "/v1/users"

	signInPath                 = "/sign-in"
	signUpPath                 = "/sign-up"
	sendConfirmationEmailPath  = "/send-confirmation-email"
	confirmEmailPath           = "/confirm-email"
	sendResetPasswordEmailPath = "/send-reset-password-email"
	resetPasswordPath          = "/reset-password"

	ownUserPath        = "/me"
	updatePasswordPath = "/password"
)

/* Send confirmation email */

const (
	sendConfirmationEmailFirstNameJoker = "<FIRSTNAME>"
	sendConfirmationEmailTokenJoker     = "<TOKEN>"
	sendConfirmationEmailLink           = "https://cine-circle-api/send-email-confirmation"
	sendConfirmationEmailSubject        = "[Cine-circle] Confirmation de votre adresse mail"
	sendConfirmationEmailBody           = "Bonjour " + sendConfirmationEmailFirstNameJoker + ",\n" +
		"Veuillez confirmer votre adresse mail en cliquant sur le lien suivant : \n" +
		sendConfirmationEmailLink + "/" + sendConfirmationEmailTokenJoker + "\n" +
		"Merci et à bientôt."
)

var (
	sendConfirmationEmailTags = []string{"send_confirmation_email"}
)

/* Send reset password email confirmation */

const (
	sendResetPasswordEmailFirstNameJoker = "<FIRSTNAME>"
	sendResetPasswordTokenJoker          = "<TOKEN>"
	sendResetPasswordEmailLink           = "https://cine-circle-api/send-reset-password-email"
	sendResetPasswordEmailSubject        = "[Cine-circle] Réinitilisation de votre mot de passe"
	sendResetPasswordEmailBody           = "Bonjour " + sendResetPasswordEmailFirstNameJoker + ",\n" +
		"Veuillez confirmer votre adresse mail en cliquant sur le lien suivant : \n" +
		sendResetPasswordEmailLink + "/" + sendResetPasswordTokenJoker + "\n" +
		"Merci et à bientôt."
)

var (
	sendResetPasswordEmailTags = []string{"send_reset_password_email"}
)
