package model

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	RatingBoundMin = 0
	RatingBoundMax = 10
)

var (
	RatingOver = "/" + fmt.Sprintf("%d", RatingBoundMax)
)

type UserRating struct {
	gorm.Model
	UserID uint `json:"userId" gorm:"index:idx_movie_user,unique"`
	MovieId string `json:"movieId" gorm:"index:idx_movie_user,unique"`
	Rating float64 `json:"rating"`
	Comment string `json:"comment"`
}