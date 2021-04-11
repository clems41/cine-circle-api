package watchlistDom

import "cine-circle/internal/domain"

type Creation struct {
	MovieID 	string 				`json:"movieId"`
	UserID 		domain.IDType 		`json:"userId"`
}

type Delete struct {
	MovieID 	string 				`json:"movieId"`
	UserID 		domain.IDType 		`json:"userId"`
}
