package repository

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/searchUtils"
)

type UserSearchForm struct {
	searchUtils.PaginationRequest
	searchUtils.SortingRequest
	Keyword string
}

type UserSearchView struct {
	searchUtils.Page
	Users []model.User
}

type User interface {
	GetFromLogin(login string) (user model.User, ok bool, err error)
	Get(userId uint) (user model.User, ok bool, err error)
	Save(user *model.User) (err error)
	Delete(userId uint) (ok bool, err error)
	Search(form UserSearchForm) (view UserSearchView, err error)
	UsernameAlreadyExists(username string) (exists bool, err error)
	EmailAlreadyExists(email string) (exists bool, err error)
}
