package watchlistDom

import (
	"cine-circle/internal/utils"
	"time"
)

type Creation struct {
	MovieID uint `json:"-"`
	UserID  uint `json:"-"`
}

type Delete struct {
	MovieID uint `json:"-"`
	UserID  uint `json:"-"`
}

type Check struct {
	MovieID uint `json:"-"`
	UserID  uint `json:"-"`
}

type List struct {
	utils.Page
	Movies []MovieView `json:"movies"`
}

type Filters struct {
	utils.PaginationRequest
	UserID uint `json:"-"`
}

type MovieView struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	PosterPath  string    `json:"posterPath"`
	ReleaseDate time.Time `json:"releaseDate"`
}
