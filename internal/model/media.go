package model

import (
	"github.com/lib/pq"
	"time"
)

type MediaType string

const (
	MovieMediaType = "movie"
	TvMediaType    = "tv"
)

type Media struct {
	ID uint `gorm:"primarykey" json:"id"`
	Metadata
	Title             string         `json:"title"`
	MediaProviderId   string         `json:"mediaProviderId" gorm:"check:media_provider_id <> ''"`
	MediaProviderName string         `json:"mediaProviderName" gorm:"check:media_provider_name <> ''"`
	Completed         bool           `json:"completed"`
	BackdropUrl       string         `json:"backdropUrl"`
	Genres            pq.StringArray `json:"genres" gorm:"type:varchar(64)[]"`
	Language          string         `json:"language"`
	OriginalTitle     string         `json:"originalTitle"`
	Overview          string         `json:"overview"`
	PosterUrl         string         `json:"posterUrl"`
	ReleaseDate       time.Time      `json:"releaseDate"`
	Runtime           int            `json:"runtime"`
	MediaType         MediaType      `json:"media_type"`
}
