package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName 	string 		`json:"FullName"`
	Ratings 	[]UserRating `json:"Ratings"`
}