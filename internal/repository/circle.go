package repository

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/searchUtils"
)

type CircleSearchForm struct {
	searchUtils.PaginationRequest
	UserId uint
}

type CircleSearchView struct {
	searchUtils.Page
	Circles []model.Circle
}

type Circle interface {
	Get(circleId uint) (circle model.Circle, ok bool, err error)
	Save(circle *model.Circle) (err error)
	Delete(circleId uint) (err error)
	Search(form CircleSearchForm) (view CircleSearchView, err error)
	AddUser(userId uint, circle *model.Circle) (err error)
	DeleteUser(userId uint, circle *model.Circle) (err error)
}
