package circleDom

import (
	"cine-circle-api/pkg/utils/searchUtils"
)

/* Common */

type CommonForm struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	UserId      uint   `json:"-"` // Get it from token
}

type UserView struct {
	Id        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
}

type CommonView struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Users       []UserView `json:"users"`
}

/* Create */

type CreateForm struct {
	CommonForm
}

type CreateView struct {
	CommonView
}

/* Add user */

type AddUserForm struct {
	UserIdToAdd uint `json:"-"` // Get it from path parameter
	UserId      uint `json:"-"` // Get it from token
	CircleId    uint `json:"-"` // Get it from path parameter
}

type AddUserView struct {
	CommonView
}

/* Delete user */

type DeleteUserForm struct {
	UserIdToDelete uint `json:"-"` // Get it from path parameter
	UserId         uint `json:"-"` // Get it from token
	CircleId       uint `json:"-"` // Get it from path parameter
}

type DeleteUserView struct {
	CommonView
}

/* Update */

type UpdateForm struct {
	CircleId uint `json:"-"` // Get it from path parameter
	CommonForm
}

type UpdateView struct {
	CommonView
}

/* Delete */

type DeleteForm struct {
	UserId   uint `json:"-"` // Get it from token
	CircleId uint `json:"-"` // Get it from path parameter
}

/* Get */

type GetForm struct {
	UserId   uint `json:"-"` // Get it from token
	CircleId uint `json:"-"` // Get it from path parameter
}

type GetView struct {
	CommonView
}

/* Search */

type SearchForm struct {
	searchUtils.PaginationRequest
	UserId uint `json:"-"` // Get it from token
}

type SearchView struct {
	searchUtils.Page
	Circles []CommonView `json:"circles"`
}
