package model

type User struct {
	ID uint `json:"id" gorm:"primarykey"`
	Metadata
	Username       string `json:"username" gorm:"uniqueIndex;check:username <> ''"`
	HashedPassword string `json:"hashedPassword" gorm:"check:hashed_password <> ''"`
	LastName       string `json:"lastName" gorm:"index;check:last_name <> ''"`
	FirstName      string `json:"firstName" gorm:"index;check:first_name <> ''"`
	Email          string `json:"email" gorm:"uniqueIndex;check:email <> ''"`
	Active         bool   `json:"active"`
	EmailToken     string `json:"emailToken"`
	PasswordToken  string `json:"passwordToken"`
	EmailConfirmed bool   `json:"emailConfirmed"`
}
