package movieDom

import (
	"cine-circle/internal/utils"
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

type SearchView struct {
	utils.Page
	Results []ItemView `json:"results"`
}

type ItemView struct {
	ID           uint   `json:"id"`
	MediaType    string `json:"mediaType"`
	Name         string `json:"name"`
	OriginalName string `json:"originalName"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"posterPath"`
}

type Filters struct {
	utils.PaginationRequest
	Query string
}

func (filters Filters) Valid() error {
	if filters.Query == "" {
		return errEmptyQuery
	}
	return nil
}
