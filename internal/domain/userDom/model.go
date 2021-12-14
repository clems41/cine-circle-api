package userDom

import (
	"time"
)

/* Common */

type CommonForm struct {
	FirstName string `json:"firstName" validate:"required,alpha"`             // Obligatoire et doit contenir uniquement des lettres
	LastName  string `json:"lastName" validate:"required,alpha"`              // Obligatoire et doit contenir uniquement des lettres
	Email     string `json:"email" validate:"required,email"`                 // Obligatoire et doit être sous la forme d'un email
	Username  string `json:"username" validate:"required,alphanum,lowercase"` // Obligatoire et doit contenir uniquement des lettres en minuscule et des chiffres
}

type CommonView struct {
	Id             uint   `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	EmailConfirmed bool   `json:"emailConfirmed"`
}

/* SignIn */

type SignInForm struct {
	Password string `json:"-"` // récupéré depuis le header Authorization (Basic authentication)
	Login    string `json:"-"` // récupéré depuis le header Authorization (Basic authentication)
}

type Token struct {
	ExpirationDate time.Time `json:"expirationDate"`
	TokenString    string    `json:"tokenString"`
}

type SignInView struct {
	CommonView
	Token Token `json:"token"`
}

/* SignUp */

type SignUpForm struct {
	CommonForm
	Password string `json:"password" validate:"required,min=8"` // Obligatoire et doit contenir au moins 8 caractères
}

type SignUpView struct {
	CommonView
}

/* Send email confirmation */

type SendEmailConfirmationForm struct {
	UserId uint `json:"-"` // Champ récupéré depuis le token
}

type ConfirmEmailForm struct {
	UserId     uint   `json:"-"` // Champ récupéré depuis le token
	EmailToken string `json:"emailToken"`
}

/* Reset password */

type SendResetPasswordEmailForm struct {
	Login string `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type ResetPasswordForm struct {
	PasswordToken string `json:"passwordToken" validate:"required"`
	Login         string `json:"login" validate:"required"`
	NewPassword   string `json:"newPassword" validate:"required,min=8"`
}

/* Update password */

type UpdatePasswordForm struct {
	UserId      uint   `json:"-"` // Champ récupéré depuis le path parameter de la route
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

/* Update info */

type UpdateForm struct {
	UserId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
	CommonForm
}

type UpdateView struct {
	CommonView
}

/* Get own info */

type GetOwnInfoForm struct {
	UserId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type GetOwnInfoView struct {
	CommonView
}

/* Delete */

type DeleteForm struct {
	UserId   uint   `json:"-"` // Champ récupéré depuis le path parameter de la route
	Password string `json:"password"`
}
