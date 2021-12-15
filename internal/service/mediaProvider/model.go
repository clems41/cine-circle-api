package mediaProvider

import (
	"time"
)

type Language string

type MediaType string

/* Get media */

type MediaView struct {
	Id            string
	Title         string
	BackdropUrl   string
	Genres        []string
	Language      string
	OriginalTitle string
	Overview      string
	PosterUrl     string
	ReleaseDate   time.Time
	Runtime       time.Duration
}

type MediaForm struct {
	Id       string
	Language Language
	Type     MediaType
}

/* Search media */

type MediaShortView struct {
	Id            string
	Title         string
	Language      string
	OriginalTitle string
	PosterUrl     string
}

type SearchForm struct {
	Page     int
	Language Language
	Keyword  string
}

type SearchView struct {
	NumberOfItems int
	NumberOfPages int
	CurrentPage   int
	Result        []MediaShortView
}
