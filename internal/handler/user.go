package handler

import (
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func NewUserHandler() *restful.WebService {
	wsUser := &restful.WebService{}
	wsUser.Path("/v1/users")

	wsUser.Route(wsUser.PUT("/{userId}").
		Doc("Update existing user").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(UpdateUser))

	wsUser.Route(wsUser.DELETE("/{userId}").
		Doc("Delete existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(DeleteUser))

	wsUser.Route(wsUser.GET("/").
		Doc("Search for user(s)").
		Param(wsUser.QueryParameter("username", "search user by username").DataType("string")).
		Param(wsUser.QueryParameter("email", "search user by email").DataType("string")).
		Param(wsUser.QueryParameter("fullname", "search user by fullname").DataType("string")).
		Writes([]model.User{}).
		Returns(200, "OK", []model.User{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(SearchUsers))

	wsUser.Route(wsUser.GET("/{userId}").
		Param(wsUser.PathParameter("userId", "ID of sought user").DataType("int")).
		Doc("Get user info from ID").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetUser))

	wsUser.Route(wsUser.GET("/me").
		Doc("Get user info from token").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetOwnUserInfo))

	wsUser.Route(wsUser.GET("/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUser.PathParameter("username", "username of sought user").DataType("string")).
		Writes("").
		Returns(200, "OK", "").
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(false)).
		To(UsernameExists))

	wsUser.Route(wsUser.GET("/{userId}/movies").
		Doc("Get all movies that user had rated").
		Param(wsUser.PathParameter("userId", "username of sought user").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetMoviesByUser))

	return wsUser
}

func CreateUser(req *restful.Request, res *restful.Response) {
	var user model.User

	err := req.ReadEntity(&user)
	if err != nil {
		res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
			typedErrors.ErrApiUnprocessableEntity.CodeError())
		return
	}

	if service.UserExists("username = ?", user.Username) {
		res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
			typedErrors.ErrApiUnprocessableEntity.CodeError())
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
			res.WriteHeaderAndEntity(typedErrors.ErrRepositoryResourceNotFound.HttpCode(), typedErrors.ErrRepositoryResourceNotFound.CodeError())
			return
		}
		err := req.ReadEntity(&user)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
				typedErrors.ErrApiUnprocessableEntity.CodeError())
			return
		}

		var err2 typedErrors.CustomError
		err2, newUser = service.CreateOrUpdateUser(user, "id = ?", userIdStr)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	}
	res.WriteHeaderAndEntity(http.StatusOK, newUser)
}

func DeleteUser(req *restful.Request, res *restful.Response) {
	userIdStr := req.PathParameter("userId")
	if userIdStr != "" {
		if !service.UserExists("id = ?", userIdStr) {
			res.WriteHeaderAndEntity(typedErrors.ErrRepositoryResourceNotFound.HttpCode(), typedErrors.ErrRepositoryResourceNotFound.CodeError())
			return
		}

		var err2 typedErrors.CustomError
		err2 = service.DeleteUser("id = ?", userIdStr)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func GetUser(req *restful.Request, res *restful.Response) {
	userId := req.PathParameter("userId")
	var user model.User
	if userId != "" {
		var err typedErrors.CustomError
		err, user = service.GetUser("id = ?", userId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func GetOwnUserInfo(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
	var user model.User
	if username != "" {
		var err typedErrors.CustomError
		err, user = service.GetUser("username = ?", username)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func SearchUsers(req *restful.Request, res *restful.Response) {
	username := req.QueryParameter("username")
	email := req.QueryParameter("email")
	fullname := req.QueryParameter("fullname")
	var users []model.User
	var err typedErrors.CustomError
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
		res.WriteHeaderAndEntity(http.StatusOK, "true")
	} else {
		res.WriteHeaderAndEntity(http.StatusOK, "false")
	}
}

func GetMoviesByUser(req *restful.Request, res *restful.Response) {
	userId := req.PathParameter("userId")
	var movies []model.Movie
	if userId != "" {
		var err typedErrors.CustomError
		err, movies = service.GetMoviesByUser("id = ?", userId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}
