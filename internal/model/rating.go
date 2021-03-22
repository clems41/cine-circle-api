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

type Rating struct {
	gorm.Model
	UserID uint `json:"UserID" gorm:"index:idx_movie_user,unique"`
	MovieID string `json:"MovieID" gorm:"index:idx_movie_user,unique"`
	Source string `json:"Source"`
	Value float64 `json:"Value"`
	Comment string `json:"Comment"`
	Username string `json:"Username"`
}