package model

import (
	"cine-circle/internal/typedErrors"
	"strings"
)

type User struct {
	GormModel
	FullName 	string 		`json:"fullName"`
	Username 	string 		`json:"username" gorm:"uniqueIndex"`
	Email	 	string 		`json:"email"`
	Hash 		string 		`json:"-"`
	Password	string 		`json:"password" gorm:"-"`
}

func (user User) IsValid() typedErrors.CustomError {
	if user.Username == "" || user.FullName == "" || user.Email == "" || user.Password == "" {
		return typedErrors.ErrServiceMissingMandatoryFields
	}
	if strings.Contains(user.Username, " ") {
		return typedErrors.ErrServiceBadFormatMandatoryFields
	}
	return typedErrors.NoErr
}