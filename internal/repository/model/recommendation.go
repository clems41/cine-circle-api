package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Recommendation struct {
	gormUtils.Metadata
	SenderId uint `gorm:"index"`
	Sender   User
	Circles  []Circle `gorm:"many2many:recommendation_circles"`
	MovieId  uint     `gorm:"index"`
	Movie    Movie
	Text     string
	Date     time.Time `gorm:"index"`
}

func MigrateRecommendation(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&Recommendation{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}
