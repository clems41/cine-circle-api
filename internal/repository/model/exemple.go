package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Exemple struct {
	gormUtils.Metadata
	// TODO add your custom fields here and their rules (cf. user model example with gorm constraints)
}

func MigrateExemple(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&Exemple{})
	if err != nil {
		return errors.WithStack(err)
	}

	return
}
