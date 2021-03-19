package model

import "gorm.io/gorm"

type Circle struct {
	gorm.Model
	Users []User `gorm:"many2many:user_circle;"`
	Name string `json:"Name"`
	Description string `json:"Description"`
}
