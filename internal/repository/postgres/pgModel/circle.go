package pgModel

import (
	"cine-circle-api/pkg/sql/gormUtils"
)

type Circle struct {
	gormUtils.Metadata
	Name        string `gorm:"index;check:name <> ''"`
	Description string
	Users       []User `gorm:"many2many:circle_users"`
}
