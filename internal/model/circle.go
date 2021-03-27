package model

type Circle struct {
	GormModel
	Users []User `gorm:"many2many:user_circle;"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type UserCircle struct {
	CircleID uint
	UserID uint
}

func (c *Circle) IsValid() CustomError {
	return NoErr
}