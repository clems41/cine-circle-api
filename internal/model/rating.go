package model

import (
	"fmt"
)

const (
	RatingBoundMin = 0
	RatingBoundMax = 10
)

var (
	RatingOver = "/" + fmt.Sprintf("%d", RatingBoundMax)
)

type Rating struct {
	GormModel
	UserID uint `json:"userId" gorm:"index:idx_movie_user,unique"`
	MovieID string `json:"movieId" gorm:"index:idx_movie_user,unique"`
	Source string `json:"source"`
	Value float64 `json:"value"`
	Comment string `json:"comment"`
	Username string `json:"username"`
}