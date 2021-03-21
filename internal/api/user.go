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

	if service.UserExists("username = ?", user.Username) {
		res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
			model.ErrInternalApiUnprocessableEntity.CodeError())
		return
	}

	err2, newUser := service.CreateOrUpdateUser(user, "username = ?", user.Username)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, newUser)
}

func UpdateUser(req *restful.Request, res *restful.Response) {
	userIdStr := req.PathParameter("userId")
	var user, newUser model.User
	if userIdStr != "" {
		if !service.UserExists("id = ?", userIdStr) {
			res.WriteHeaderAndEntity(model.ErrInternalDatabaseResourceNotFound.HttpCode(), model.ErrInternalDatabaseResourceNotFound.CodeError())
			return
		}
		err := req.ReadEntity(&user)
		if err != nil {
			res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
				model.ErrInternalApiUnprocessableEntity.CodeError())
			return
		}

		var err2 model.CustomError
		err2, newUser = service.CreateOrUpdateUser(user, "id = ?", userIdStr)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
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

func SearchUsers(req *restful.Request, res *restful.Response) {
	username := req.QueryParameter("username")
	email := req.QueryParameter("email")
	fullname := req.QueryParameter("fullname")
	var users []model.User
	var err model.CustomError
	err, users = service.SearchUsers(username, fullname, email)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, users)
}

func UsernameExists(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")
	if service.UserExists("username = ?", username) {
		res.WriteHeaderAndEntity(http.StatusFound, "true")
	} else {
		res.WriteHeaderAndEntity(http.StatusNotFound, "false")
	}
}

func GetMoviesByUser(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")
	var movies []model.Movie
	if username != "" {
		var err model.CustomError
		err, movies = service.GetMoviesByUser(username)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}
