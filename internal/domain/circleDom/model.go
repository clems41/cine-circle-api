package circleDom

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
)

type Creation struct {
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	UsersID 		[]domain.IDType	`json:"usersId"`
}

type Update struct {
	CircleID		domain.IDType 	`json:"id"`
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	UsersID 		[]domain.IDType	`json:"usersId"`
}

type Delete struct {
	CircleID		domain.IDType 	`json:"id"`
}

type Get struct {
	CircleID		domain.IDType 	`json:"id"`
}

type Result struct {
	CircleID		domain.IDType 	`json:"id"`
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	Users 			[]userDom.Result		`json:"users"`
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