package recommendationDom

import (
	"cine-circle/internal/utils"
	"time"
)

type Creation struct {
	SenderID  uint   `json:"-"`
	MovieID   uint   `json:"movieId"`
	Comment   string `json:"comment"`
	CircleIDs []uint `json:"circleIds"`
	UserIDs   []uint `json:"userIds"`
}

type UserView struct {
	UserID      uint   `json:"userId"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type CircleView struct {
	CircleID    uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []UserView `json:"users"`
}

type MovieView struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	PosterPath  string    `json:"posterPath"`
	ReleaseDate time.Time `json:"releaseDate"`
}

type RecommendationView struct {
	ID                 uint         `json:"id"`
	Date               time.Time    `json:"date"`
	Sender             UserView     `json:"sender"`
	Movie              MovieView    `json:"movie"`
	Comment            string       `json:"comment"`
	Circles            []CircleView `json:"circles"`
	Users              []UserView   `json:"users"`
	RecommendationType string       `json:"type"`
}

type ViewList struct {
	utils.Page
	Recommendations []RecommendationView `json:"recommendations"`
}

type UserList struct {
	utils.Page
	Users []UserView `json:"users"`
}

type UsersFilters struct {
	utils.PaginationRequest
	UserID uint `json:"-"`
}

type Filters struct {
	utils.PaginationRequest
	utils.SortingRequest
	RecommendationType string `json:"type"`
	UserID             uint   `json:"-"`
	MovieID            uint   `json:"-"`
	CircleID           uint   `json:"-"`
}

func (c Creation) Valid() error {
	if c.SenderID == 0 {
		return errSenderIDNull
	}
	if c.MovieID == 0 {
		return errMovieIDNull
	}
	if c.Comment == "" {
		return errCommentEmpty
	}
	if len(c.CircleIDs) == 0 && len(c.UserIDs) == 0 {
		return errMissingRecipient
	}
	return nil
}

func (f Filters) Valid() error {
	if !utils.SliceContainsStr(acceptedTypeOfRecommendation, f.RecommendationType) {
		return errRecommendationTypeIncorrect
	}
	if !utils.SliceContainsStr(acceptedFieldsForSorting, f.Field) {
		return errSortingFieldIncorrect
	}
	return nil
}
