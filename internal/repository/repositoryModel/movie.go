package repositoryModel

import (
	"github.com/lib/pq"
	"time"
)

type Movie struct {
	Metadata
	Title            string `gorm:"not null;default:null"`
	ImdbId           string
	BackdropPath     string
	PosterPath       string
	Genres           pq.StringArray `gorm:"type:text[]"`
	OriginalLanguage string
	OriginalTitle    string
	Overview         string
	ReleaseDate      time.Time
	Runtime          int
}
