package model

import "cine-circle/internal/typedErrors"

type Circle struct {
	GormModel
	Users []User `gorm:"many2many:user_circle;" json:"users"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type UserCircle struct {
	CircleID uint
	UserID uint
}

func (c *Circle) IsValid() typedErrors.CustomError {
	return typedErrors.NoErr
}