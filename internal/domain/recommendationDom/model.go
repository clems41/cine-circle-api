package recommendationDom

import (
	"cine-circle-api/pkg/utils/searchUtils"
	"time"
)

/* Common */

type CommonForm struct {
	SenderId   uint   `json:"-"` // Get it from token
	CirclesIds []uint `json:"circlesIds" validate:"min=1"`
	MediaId    uint   `json:"mediaId" validate:"ne=0"`
	Text       string `json:"text" validate:"required"`
}

type CommonView struct {
	Id      uint         `json:"id"`
	Sender  UserView     `json:"sender"`
	Circles []CircleView `json:"circles"`
	Movie   MovieView    `json:"movie"`
	Text    string       `json:"text"`
	Date    time.Time    `json:"date"`
	Type    string       `json:"type"`
}

type UserView struct {
	Id        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
}

type CircleView struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []UserView `json:"users"`
}

type MovieView struct {
	Id            uint      `json:"id"`
	Title         string    `json:"title"`
	BackdropUrl   string    `json:"backdropUrl"`
	Genres        []string  `json:"genres"`
	Language      string    `json:"language"`
	OriginalTitle string    `json:"originalTitle"`
	Overview      string    `json:"overview"`
	PosterUrl     string    `json:"posterUrl"`
	ReleaseDate   time.Time `json:"releaseDate"`
	Runtime       int       `json:"runtime"`
}

/* Send */

type SendForm struct {
	CommonForm
}

type SendView struct {
	CommonView
}

/* Search */

type SearchForm struct {
	searchUtils.PaginationRequest
	MovieId int    `json:"mediaId"`
	Type    string `json:"type" validate:"oneof=all sent received"`
	UserId  uint   `json:"-"` // Get it from token
}

type SearchView struct {
	searchUtils.Page
	Recommendations []CommonView `json:"circles"`
}
