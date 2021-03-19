package api

import (
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func CreateUser(req *restful.Request, res *restful.Response) {
	var user model.User

	err := req.ReadEntity(&user)
	if err != nil {
		res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
			model.ErrInternalApiUnprocessableEntity.CodeError())
		return
	}

	err2 := service.CreateUser(user.FullName)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, "")
}

func GetUser(req *restful.Request, res *restful.Response) {
	userId := req.QueryParameter("id")
	var user model.User
	if userId != "" {
		var err model.CustomError
		err, user = service.GetUser(userId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, user)
}
