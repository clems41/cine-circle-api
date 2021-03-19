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

	err2, newUser := service.CreateUser(user.Username, user.FullName, user.Email)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, newUser)
}

func GetUser(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")
	var user model.User
	if username != "" {
		var err model.CustomError
		err, user = service.GetUser(username)
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

func UsernameExists(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")
	if service.UsernameAlreadyExists(username) {
		res.WriteHeaderAndEntity(http.StatusFound, "true")
	} else {
		res.WriteHeaderAndEntity(http.StatusNotFound, "false")
	}
}
