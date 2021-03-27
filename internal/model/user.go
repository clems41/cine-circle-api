package model

import (
	"gorm.io/gorm"
	"strings"
)

type User struct {
	gorm.Model
	FullName 	string 		`json:"FullName"`
	Username 	string 		`json:"Username" gorm:"uniqueIndex"`
	Email	 	string 		`json:"Email"`
	Hash 		string 		`json:"-"`
	Password	string 		`json:"Password" gorm:"-"`
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