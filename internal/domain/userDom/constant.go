package userDom

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
