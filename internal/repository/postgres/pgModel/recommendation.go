package pgModel

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"time"
)

type Recommendation struct {
	gormUtils.Metadata
	SenderId uint `gorm:"index"`
	Sender   User
	Circles  []Circle `gorm:"many2many:recommendation_circles"`
	MediaID  uint     `gorm:"index"`
	Media    Media
	Text     string
	Date     time.Time `gorm:"index"`
}
