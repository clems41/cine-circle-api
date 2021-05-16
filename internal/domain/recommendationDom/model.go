package recommendationDom

import (
	"time"
)

type Creation struct {
	MovieID   string `json:"movieId"`
	UsersId   []uint `json:"usersId"`
	CirclesId []uint `json:"circlesId"`
	Comment   string `json:"comment"`
}

type UserView struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
}

type CircleView struct {
	CircleID uint   `json:"id"`
	Name     string `json:"name"'`
}

type Result struct {
	AddedAt time.Time    `json:"addedAt"`
	MovieID string       `json:"movieId"`
	Users   []UserView   `json:"users"`
	Circles []CircleView `json:"circles"`
	Comment string       `json:"comment"`
}
