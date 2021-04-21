package handler

import (
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"github.com/emicklei/go-restful"
	"net/http"
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

	wsUser.Route(wsUser.PUT("/{userId}").
		Doc("Update existing user").
		Writes(userDom.Update{}).
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Update))

	wsUser.Route(wsUser.PUT("/{userId}/password").
		Doc("Update existing user's password").
		Writes(userDom.UpdatePassword{}).
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.UpdatePassword))

	wsUser.Route(wsUser.DELETE("/{userId}").
		Doc("Delete existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Delete))

	wsUser.Route(wsUser.GET("/{userId}").
		Doc("Get existing user").
		Writes("").
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Get))

	/*	wsUser.Route(wsUser.GET("/").
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

func (api userHandler) Update(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	var update userDom.Update
	err = req.ReadEntity(&update)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	update.UserID = userID

	user, err := api.service.Update(update)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (api userHandler) UpdatePassword(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	var updatePassword userDom.UpdatePassword
	err = req.ReadEntity(&updatePassword)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	updatePassword.UserID = userID

	user, err := api.service.UpdatePassword(updatePassword)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (api userHandler) Delete(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	err = api.service.Delete(userDom.Delete{UserID: userID})
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (api userHandler) Get(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	user, err := api.service.Get(userDom.Get{UserID: userID})
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}