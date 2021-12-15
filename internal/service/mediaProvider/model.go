package mediaProvider

import (
	"time"
)

/* Get movie */

type MovieView struct {
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

type MovieForm struct {
	Id string
}

/* Search movie */

type MovieShortView struct {
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

type SearchView struct {
	NumberOfItems int
	NumberOfPages int
	CurrentPage   int
	Result        []MovieShortView
}
