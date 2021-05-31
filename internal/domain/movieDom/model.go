package movieDom

import (
	"time"
)

type View struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	ImdbId           string    `json:"imdbId"`
	BackdropPath     string    `json:"backdropPath"`
	PosterPath       string    `json:"posterPath"`
	Genres           []string  `json:"genres"`
	OriginalLanguage string    `json:"originalLanguage"`
	OriginalTitle    string    `json:"originalTitle"`
	Overview         string    `json:"overview"`
	ReleaseDate      time.Time `json:"releaseDate"`
	Runtime          int       `json:"runtime"`
	Trailer          string    `json:"trailer"`
}

type QueryParameter struct {
	Key   string
	Value string
}
