package mediaDom

import (
	"cine-circle-api/pkg/utils/searchUtils"
	"time"
)

type GetForm struct {
	MediaId uint
}

type GetView struct {
	Id            uint          `json:"id"`
	Title         string        `json:"title"`
	BackdropUrl   string        `json:"backdropUrl"`
	Genres        []string      `json:"genres"`
	Language      string        `json:"language"`
	OriginalTitle string        `json:"originalTitle"`
	Overview      string        `json:"overview"`
	PosterUrl     string        `json:"posterUrl"`
	ReleaseDate   time.Time     `json:"releaseDate"`
	Runtime       time.Duration `json:"runtime"`
}

type SearchForm struct {
	searchUtils.PaginationRequest
	Keyword string `json:"keyword"`
}

type SearchView struct {
	searchUtils.Page
	Result []ResultView `json:"result"`
}

type ResultView struct {
	Id            uint   `json:"id"`
	Title         string `json:"title"`
	Language      string `json:"language"`
	OriginalTitle string `json:"originalTitle"`
	PosterUrl     string `json:"posterUrl"`
}
