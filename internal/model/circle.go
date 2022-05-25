package model

type Circle struct {
	ID uint `json:"id" gorm:"primarykey"`
	Metadata
	Name        string `json:"name" gorm:"index;check:name <> ''"`
	Description string `json:"description"`
	Users       []User `json:"users" gorm:"many2many:circle_users"`
}
