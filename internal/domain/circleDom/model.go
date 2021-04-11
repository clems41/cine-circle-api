package circleDom

import (
	"cine-circle/internal/domain"
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

type Result struct {
	CircleID		domain.IDType 	`json:"id"`
	Name 			string 			`json:"name"`
	Description 	string 			`json:"description"`
	Users 			[]UserView		`json:"users"`
}

type UserView struct {
	UserID 			domain.IDType 	`json:"id"`
	Username 		string 			`json:"username"`
}