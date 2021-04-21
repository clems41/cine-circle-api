package userDom

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/typedErrors"
)

type Creation struct {
	Username 		string 			`json:"username"`
	DisplayName 	string 			`json:"displayName"`
	Password 		string			`json:"password"`
	Email 			string 			`json:"email"`
}

type Update struct {
	UserID 			domain.IDType 	`json:"-"`
	DisplayName 	string 			`json:"displayName"`
	Email 			string 			`json:"email"`
}

type Get struct {
	UserID 			domain.IDType 	`json:"-"`
	Username 		string 			`json:"username"`
	Email 			string 			`json:"email"`
}

type UpdatePassword struct {
	UserID 				domain.IDType 	`json:"-"`
	OldPassword 		string			`json:"oldPassword"`
	NewPassword 		string			`json:"newPassword"`
	NewHashedPassword 	string			`json:"-"`
}

type Delete struct {
	UserID 			domain.IDType 	`json:"-"`
}

type Result struct {
	UserID 			domain.IDType 	`json:"id"`
	Username 		string 			`json:"username"`
	DisplayName 	string 			`json:"displayName"`
	Email 			string 			`json:"email"`
}

var (
	errValidPassword = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Password is empty")
	errValidEmail = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Email is empty")
	errValidUsername = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Username is empty")
	errValidDisplayName = typedErrors.NewServiceMissingMandatoryFieldsErrorf("DisplayName is empty")
	errValidID = typedErrors.NewServiceMissingMandatoryFieldsErrorf("UserID is empty")
	errValidOldPassword = typedErrors.NewServiceMissingMandatoryFieldsErrorf("OldPassword is empty")
	errValidNewPassword = typedErrors.NewServiceMissingMandatoryFieldsErrorf("NewPassword is empty")
)

func (c Creation) Valid() (err error) {
	if c.Password == "" {
		err = errValidPassword
	}
	if c.Email == "" {
		err = errValidEmail
	}
	if c.Username == "" {
		err = errValidUsername
	}
	if c.DisplayName == "" {
		err = errValidDisplayName
	}
	return
}

func (u Update) Valid() (err error) {
	if u.UserID == 0 {
		err = errValidID
	}
	if u.Email == "" {
		err = errValidEmail
	}
	if u.DisplayName == "" {
		err = errValidDisplayName
	}
	return
}

func (d Delete) Valid() (err error) {
	if d.UserID == 0 {
		err = errValidID
	}
	return
}

func (up UpdatePassword) Valid() (err error) {
	if up.UserID == 0 {
		err = errValidID
	}
	if up.OldPassword == "" {
		err = errValidOldPassword
	}
	if up.NewPassword == "" {
		err = errValidNewPassword
	}
	return
}
