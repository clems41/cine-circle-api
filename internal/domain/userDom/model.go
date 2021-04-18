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
	UserID 			domain.IDType 	`json:"id"`
	DisplayName 	string 			`json:"displayName"`
	Email 			string 			`json:"email"`
}

type Get struct {
	UserID 			domain.IDType 	`json:"id"`
	Username 		string 			`json:"username"`
	Email 			string 			`json:"email"`
}

type UpdatePassword struct {
	OldPassword 	string			`json:"oldPassword"`
	NewPassword 	string			`json:"newPassword"`
}

type Delete struct {
	UserID 			domain.IDType 	`json:"id"`
}

type Result struct {
	UserID 			domain.IDType 	`json:"id"`
	Username 		string 			`json:"username"`
	DisplayName 	string 			`json:"displayName"`
	Email 			string 			`json:"email"`
}

func (c Creation) Valid() (err error) {
	if c.Password != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Password is empty")
	}
	if c.Email != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Email is empty")
	}
	if c.Username != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Username is empty")
	}
	if c.DisplayName != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("DisplayName is empty")
	}
	return
}

func (u Update) Valid() (err error) {
	if u.Email != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Email is empty")
	}
	if u.DisplayName != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("DisplayName is empty")
	}
	return
}

func (up UpdatePassword) Valid() (err error) {
	if up.OldPassword != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("OldPassword is empty")
	}
	if up.NewPassword != "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("NewPassword is empty")
	}
	return
}
