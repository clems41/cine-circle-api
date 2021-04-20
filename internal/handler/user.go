package handler

import (
	"cine-circle/internal/domain/userDom"
	"github.com/emicklei/go-restful"
)

type userHandler struct {
	service userDom.Service
}

func NewUserHandler(svc userDom.Service) *userHandler {
	return &userHandler{
		service:    svc,
	}
}

func (api userHandler) WebService() *restful.WebService {
	wsUser := &restful.WebService{}
	wsUser.Path("/v1/users")

/*	wsUser.Route(wsUser.PUT("/{userId}").
		Doc("Update existing user").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(UpdateUser))

	wsUser.Route(wsUser.DELETE("/{userId}").
		Doc("Delete existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(DeleteUser))

	wsUser.Route(wsUser.GET("/").
		Doc("Search for user(s)").
		Param(wsUser.QueryParameter("username", "search user by username").DataType("string")).
		Param(wsUser.QueryParameter("email", "search user by email").DataType("string")).
		Param(wsUser.QueryParameter("fullname", "search user by fullname").DataType("string")).
		Writes([]model.User{}).
		Returns(200, "OK", []model.User{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(SearchUsers))

	wsUser.Route(wsUser.GET("/{userId}").
		Param(wsUser.PathParameter("userId", "ID of sought user").DataType("int")).
		Doc("Get user info from ID").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(authenticateUser(true)).
		To(Get))

	wsUser.Route(wsUser.GET("/me").
		Doc("Get user info from token").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(authenticateUser(true)).
		To(GetOwnUserInfo))

	wsUser.Route(wsUser.GET("/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUser.PathParameter("username", "username of sought user").DataType("string")).
		Writes("").
		Returns(200, "OK", "").
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(authenticateUser(false)).
		To(UsernameExists))

	wsUser.Route(wsUser.GET("/{userId}/movies").
		Doc("Get all movies that user had rated").
		Param(wsUser.PathParameter("userId", "username of sought user").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(authenticateUser(true)).
		To(GetMoviesByUser))*/

	return wsUser
}