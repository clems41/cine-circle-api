package pgModel

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/lib/pq"
	"time"
)

type Media struct {
	gormUtils.Metadata
	Title             string
	MediaProviderId   string `gorm:"check:media_provider_id <> ''"`
	MediaProviderName string `gorm:"check:media_provider_name <> ''"`
	Completed         bool
	BackdropUrl       string
	Genres            pq.StringArray `gorm:"type:varchar(64)[]"`
	Language          string
	OriginalTitle     string
	Overview          string
	PosterUrl         string
	ReleaseDate       time.Time
	Runtime           int
	MediaType         string
}
