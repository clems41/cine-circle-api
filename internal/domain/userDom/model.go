package userDom

import "cine-circle/internal/domain"

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
