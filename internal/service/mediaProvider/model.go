package mediaProvider

import (
	"cine-circle-api/internal/constant/languageConst"
	"cine-circle-api/pkg/utils/sliceUtils"
	"time"
)

type Language string

func (l Language) Valid() bool {
	return sliceUtils.SliceContainsStr(languageConst.AllowedLanguages(), string(l))
}

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
	Page    int
	Keyword string
}

type SearchView struct {
	NumberOfItems int
	NumberOfPages int
	CurrentPage   int
	Result        []MediaShortView
}
