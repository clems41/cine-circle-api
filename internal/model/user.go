package model

import (
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

func (user User) IsValid() CustomError {
	if user.Username == "" || user.FullName == "" || user.Email == "" || user.Password == "" {
		return ErrInternalServiceMissingMandatoryFields
	}
	if strings.Contains(user.Username, " ") {
		return ErrInternalServiceBadFormatMandatoryFields
	}
	return NoErr
}