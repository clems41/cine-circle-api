package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Circle struct {
	gormUtils.Metadata
	Name        string `gorm:"index;check:name <> ''"`
	Description string
	Users       []User `gorm:"many2many:circle_users"`
}

func MigrateCircle(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&Circle{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}
