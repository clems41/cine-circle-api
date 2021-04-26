package handler

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"github.com/emicklei/go-restful"
	"net/http"
)

const (
	usersBasePath = "/users"
)

type userHandler struct {
	service userDom.Service
}

func NewUserHandler(svc userDom.Service) *userHandler {
	return &userHandler{
		service:    svc,
	}
}

func (handler userHandler) WebService() *restful.WebService {
	wsUser := &restful.WebService{}

	// Route for signup or signin
	wsUser.Path("/v1")

	wsUser.Route(wsUser.POST("/signup").
		Doc("Create new user (signup)").
		Writes(userDom.Creation{}).
		Returns(201, "Created", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		To(handler.CreateUser))

	wsUser.Route(wsUser.POST("/signin").
		Doc("Generate token from username and password (basic authentication)").
		Writes(nil).
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		To(handler.GenerateToken))

	// Route for updating or getting user's info
	//wsUser.Path("/v1/users")

	wsUser.Route(wsUser.PUT(usersBasePath + "/{userId}").
		Doc("Update existing user").
		Writes(userDom.Update{}).
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.Update))

	wsUser.Route(wsUser.PUT(usersBasePath + "/{userId}/password").
		Doc("Update existing user's password").
		Writes(userDom.UpdatePassword{}).
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.UpdatePassword))

	wsUser.Route(wsUser.DELETE(usersBasePath + "/{userId}").
		Doc("Delete existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.Delete))

	wsUser.Route(wsUser.GET(usersBasePath + "/{userId}").
		Doc("Get existing user").
		Writes("").
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.Get))

	wsUser.Route(wsUser.GET(usersBasePath + "/me").
		Doc("Get user info from token").
		Writes("").
		Returns(200, "OK", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.GetOwnUserInfo))

	wsUser.Route(wsUser.GET(usersBasePath + "/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUser.PathParameter("username", "username of sought user").DataType("string")).
		Writes("").
		Returns(200, "OK", true).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.UsernameExists))

	wsUser.Route(wsUser.GET(usersBasePath + "/").
		Doc("Search for user(s)").
		Param(wsUser.QueryParameter("search", "search user using keyword (will match username, email and displayName").DataType("string")).
		Writes("").
		Returns(200, "OK", []userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(handler.SearchUsers))

	/*	wsUser.Route(wsUser.GET("/{userId}/movies").
		Doc("Get all movies that user had rated").
		Param(wsUser.PathParameter("userId", "username of sought user").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(authenticateUser(true)).
		To(GetMoviesByUser))*/

	return wsUser
}

func (handler userHandler) GenerateToken(req *restful.Request, res *restful.Response) {
	auth := req.HeaderParameter(constant.AuthenticationHeaderName)
	token, err := handler.service.GenerateTokenFromAuthenticationHeader(auth)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, token)
}

func (handler userHandler) CreateUser(req *restful.Request, res *restful.Response) {
	var userCreation userDom.Creation
	err := req.ReadEntity(&userCreation)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestError(err))
		return
	}
	user, err := handler.service.Create(userCreation)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (handler userHandler) Update(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	var update userDom.Update
	err = req.ReadEntity(&update)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestError(err))
		return
	}

	update.UserID = userID

	user, err := handler.service.Update(update)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (handler userHandler) UpdatePassword(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	var updatePassword userDom.UpdatePassword
	err = req.ReadEntity(&updatePassword)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestError(err))
		return
	}

	updatePassword.UserID = userID

	user, err := handler.service.UpdatePassword(updatePassword)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (handler userHandler) Delete(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	err = handler.service.Delete(userDom.Delete{UserID: userID})
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (handler userHandler) Get(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	user, err := handler.service.Get(userDom.Get{UserID: userID})
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (handler userHandler) GetOwnUserInfo(req *restful.Request, res *restful.Response) {
	user, err := CommonHandler.WhoAmI(req)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, user)
}

func (handler userHandler) UsernameExists(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")

	user, err := handler.service.Get(userDom.Get{Username: username})

	var exists bool
	if err != nil {
		exists = false
	} else {
		exists = user.Username == username
	}

	res.WriteHeaderAndEntity(http.StatusOK, exists)
}

func (handler userHandler) SearchUsers(req *restful.Request, res *restful.Response) {
	keyword := req.QueryParameter("search")

	users, err := handler.service.Search(userDom.Filters{Keyword: keyword})
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, users)
}