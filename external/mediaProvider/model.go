package mediaProvider

import (
	"time"
)

type Media struct {
	Id            string
	Title         string
	BackdropUrl   string
	Genres        []string
	Language      string
	OriginalTitle string
	Overview      string
	PosterUrl     string
	ReleaseDate   time.Time
	Runtime       int
}

type MediaShort struct {
	Id            string
	Title         string
	Language      string
	OriginalTitle string
	PosterUrl     string
}

type SearchForm struct {
	Page    int
	Keyword string
}
