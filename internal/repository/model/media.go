package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Media struct {
	gormUtils.Metadata
	Title string `gorm:"uniqueIndex;check:title <> ''"`
}

func MigrateMedia(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&Media{})
	if err != nil {
		return errors.WithStack(err)
	}

	return
}
