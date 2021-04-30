package userDom

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/typedErrors"
)

type Creation struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

type Update struct {
	UserID      domain.IDType `json:"-"`
	DisplayName string        `json:"displayName"`
	Email       string        `json:"email"`
}

type Get struct {
	UserID   domain.IDType `json:"-"`
	Username string        `json:"-"`
	Email    string        `json:"-"`
}

type UpdatePassword struct {
	UserID            domain.IDType `json:"-"`
	OldPassword       string        `json:"oldPassword"`
	NewPassword       string        `json:"newPassword"`
	NewHashedPassword string        `json:"-"`
}

type Delete struct {
	UserID domain.IDType `json:"-"`
}

type Result struct {
	UserID      domain.IDType `json:"id"`
	Username    string        `json:"username"`
	DisplayName string        `json:"displayName"`
}

type Filters struct {
	Keyword string
}

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
	if u.Email == "" && u.DisplayName == "" {
		err = errValidUpdate
	}
	return
}

func (d Delete) Valid() (err error) {
	if d.UserID == 0 {
		err = errValidID
	}
	return
}

func (g Get) Valid() (err error) {
	if g.UserID == 0 && g.Username == "" && g.Email == "" {
		err = errValidGet
	}
	return
}

func (f Filters) Valid() (err error) {
	if len(f.Keyword) < 3 {
		err = errValidKeyword
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
