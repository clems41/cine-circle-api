package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Movie struct {
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
}

type TvShow struct {
	gormUtils.Metadata
	Title string `gorm:"check:title <> ''"`
}

func MigrateMedia(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&Movie{})
	if err != nil {
		return errors.WithStack(err)
	}
	err = DB.Exec("ALTER TABLE movies ADD UNIQUE (media_provider_id, media_provider_name)").Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = DB.
		AutoMigrate(&TvShow{})
	if err != nil {
		return errors.WithStack(err)
	}

	return
}
