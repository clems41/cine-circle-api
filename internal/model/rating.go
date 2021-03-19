package model

import "gorm.io/gorm"

type UserRating struct {
	gorm.Model
	UserID uint `json:"userId"`
	MovieId string `json:"movieId"`
	Rating float64 `json:"rating"`
	Comment string `json:"comment"`
}