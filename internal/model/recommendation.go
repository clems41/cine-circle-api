package model

import (
	"time"
)

type Recommendation struct {
	ID      uint      `json:"id"`
	Sender  User      `json:"sender"`
	Circles []Circle  `json:"circles"`
	Media   Media     `json:"media"`
	Text    string    `json:"text"`
	Date    time.Time `json:"date"`
}
