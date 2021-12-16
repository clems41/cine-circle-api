package circleDom

import (
	"cine-circle-api/pkg/utils/searchUtils"
)

/* Common */

type CommonForm struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
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
	UserId   uint `json:"-"` // Champ récupéré depuis le path parameter de la route
	CircleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type AddUserView struct {
	CommonView
}

/* Delete user */

type DeleteUserForm struct {
	UserId   uint `json:"-"` // Champ récupéré depuis le path parameter de la route
	CircleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type DeleteUserView struct {
	CommonView
}

/* Update */

type UpdateForm struct {
	CircleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
	CommonForm
}

type UpdateView struct {
	CommonView
}

/* Delete */

type DeleteForm struct {
	CircleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

/* Get */

type GetForm struct {
	CircleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type GetView struct {
	CommonView
}

/* Search */

type SearchForm struct {
	searchUtils.PaginationRequest
	CircleName string `json:"name" validate:"required"`
}

type SearchView struct {
	searchUtils.Page
	Circles []CommonView `json:"circles"`
}
