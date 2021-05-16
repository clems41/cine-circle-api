package circleDom

import (
	"cine-circle/internal/typedErrors"
)

type Creation struct {
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	UsersID 		[]uint	`json:"usersId"`
}

type Update struct {
	CircleID		uint 	`json:"id"`
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	UsersID 		[]uint	`json:"usersId"`
}

type Delete struct {
	CircleID		uint 	`json:"id"`
}

type Get struct {
	CircleID		uint 	`json:"id"`
}

type View struct {
	CircleID		uint 	`json:"id"`
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	Users 			[]UserView		`json:"users"`
}

type UserView struct {
	UserID      uint `json:"id"`
	Username    string        `json:"username"`
	DisplayName string        `json:"displayName"`
}

func (c Creation) Valid() (err error) {
	if c.Name == "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Name should be specified")
	}
	return
}

func (u Update) Valid() (err error) {
	if u.Name == "" {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("Name should be specified")
	}
	if u.CircleID == 0 {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("ID should be specified")
	}
	return
}

func (d Delete) Valid() (err error) {
	if d.CircleID == 0 {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("ID should be specified")
	}
	return
}

func (g Get) Valid() (err error) {
	if g.CircleID == 0 {
		err = typedErrors.NewServiceMissingMandatoryFieldsErrorf("ID should be specified")
	}
	return
}
