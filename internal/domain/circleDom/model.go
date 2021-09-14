package circleDom

import (
	"cine-circle/pkg/utils"
)

type Creation struct {
	UserIDFromRequest uint   `json:"-"`
	Name              string `json:"name"`
	Description       string `json:"description"`
}

type Update struct {
	CircleID          uint   `json:"-"`
	UserIDFromRequest uint   `json:"-"`
	Name              string `json:"name"`
	Description       string `json:"description"`
}

type UpdateUser struct {
	CircleID          uint
	UserIDToUpdate    uint
	UserIDFromRequest uint
}

type Deletion struct {
	CircleID          uint `json:"-"`
	UserIDFromRequest uint `json:"-"`
}

type Get struct {
	CircleID          uint `json:"-"`
	UserIDFromRequest uint `json:"-"`
}

type View struct {
	CircleID    uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []UserView `json:"users"`
}

type UserView struct {
	UserID      uint   `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type ListView struct {
	utils.Page
	Circles []View `json:"circles"`
}

type Filters struct {
	utils.PaginationRequest
	UserID uint `json:"-"`
}

func (c Creation) Valid() (err error) {
	if c.Name == "" {
		return errNameEmpty
	}
	if c.Description == "" {
		return errDescriptionEmpty
	}
	return
}

func (u Update) Valid() (err error) {
	if u.CircleID == 0 {
		return errIdNull
	}
	if u.Name == "" || u.Description == "" {
		return errNoFieldsProvided
	}
	return
}

func (d Deletion) Valid() (err error) {
	if d.CircleID == 0 {
		return errIdNull
	}
	return
}

func (g Get) Valid() (err error) {
	if g.CircleID == 0 {
		return errIdNull
	}
	return
}
