package model

import (
	"time"
)

type Recommendation struct {
	ID uint `json:"id" gorm:"primarykey"`
	Metadata
	SenderID uint      `json:"-" gorm:"index"`
	Sender   User      `json:"sender"`
	Circles  []Circle  `json:"circles" gorm:"many2many:recommendation_circles"`
	MediaID  uint      `json:"-" gorm:"index"`
	Media    Media     `json:"media"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date" gorm:"index"`
}
