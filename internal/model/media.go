package model

import (
	"time"
)

type MediaType string

const (
	MovieMediaType = "movie"
	TvMediaType    = "tv"
)

type Media struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	BackdropUrl   string    `json:"backdrop_url"`
	Genres        []string  `json:"genres"`
	Language      string    `json:"language"`
	OriginalTitle string    `json:"original_title"`
	Overview      string    `json:"overview"`
	PosterUrl     string    `json:"poster_url"`
	ReleaseDate   time.Time `json:"release_date"`
	Runtime       int       `json:"runtime"`
	MediaType     MediaType `json:"media_type"`
}
