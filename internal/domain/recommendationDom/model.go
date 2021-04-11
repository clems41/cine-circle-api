package recommendationDom

import (
	"cine-circle/internal/domain"
	"time"
)

type Creation struct {
	MovieID 		string 				`json:"movieId"`
	UsersId 		[]domain.IDType 	`json:"usersId"`
	CirclesId 		[]domain.IDType 	`json:"circlesId"`
	Comment 		string 				`json:"comment"`
}

type UserView struct {
	UserID 			domain.IDType 		`json:"id"`
	Username 		string 				`json:"username"`
}

type CircleView struct {
	CircleID 		domain.IDType		`json:"id"`
	Name 			string 				`json:"name"'`
}

type Result struct {
	AddedAt			time.Time 			`json:"addedAt"`
	MovieID 		string 				`json:"movieId"`
	Users 			[]UserView 			`json:"users"`
	Circles 		[]CircleView 		`json:"circles"`
	Comment 		string 				`json:"comment"`
}