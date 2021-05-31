package userDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	webServicePkg "cine-circle/internal/webService"
	"cine-circle/pkg/logger"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service:    svc,
	}
}

func (ws handler) WebServices() (webServices []*restful.WebService) {
	wsAuthentication := &restful.WebService{}
	wsUser := &restful.WebService{}
	webServices = append(webServices, wsAuthentication, wsUser)

	// Route for signup or sign-in
	wsAuthentication.Path("/v1")

	wsAuthentication.Route(wsAuthentication.POST("/sign-up").
		Doc("Create new user (signup)").
		Writes(Creation{}).
		Returns(http.StatusCreated, "Created", View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		To(ws.CreateUser))

	wsAuthentication.Route(wsAuthentication.POST("/sign-in").
		Doc("Generate token from username and password (basic authentication)").
		Writes(nil).
		Returns(http.StatusOK, "OK", "token").
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		To(ws.GenerateToken))

	// Route for updating or getting user's info
	wsUser.Path("/v1/users")

	wsUser.Route(wsUser.PUT("/").
		Doc("update actual user from token").
		Writes(Update{}).
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.Update))

	wsUser.Route(wsUser.PUT("/password").
		Doc("update existing user's password").
		Writes(UpdatePassword{}).
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.UpdatePassword))

	wsUser.Route(wsUser.DELETE("/").
		Doc("Delete existing user").
		Writes(nil).
		Returns(http.StatusNoContent, "Deleted", nil).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.Delete))

	wsUser.Route(wsUser.GET("/{userId}").
		Doc("Get existing user").
		Writes(nil).
		Returns(http.StatusFound, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.Get))

	wsUser.Route(wsUser.GET("/me").
		Doc("Get user info from token").
		Writes(nil).
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.GetOwnUserInfo))

	wsUser.Route(wsUser.GET("/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUser.PathParameter("username", "username of sought user").DataType("string")).
		Writes(nil).
		Returns(http.StatusFound, "OK", true).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		To(ws.UsernameExists))

	wsUser.Route(wsUser.GET("/").
		Doc("Search for user(s)").
		Param(wsUser.QueryParameter("search", "search user using keyword (will match username, email and displayName").DataType("string")).
		Writes(nil).
		Returns(http.StatusOK, "OK", []View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(ws.SearchUsers))

	/*	wsUser.Route(wsUser.GET("/{userId}/movies").
		Doc("Get all movies that user had rated").
		Param(wsUser.PathParameter("userId", "username of sought user").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(http.StatusOK, "OK", model.MovieSearch{}).
		Returns(404, "User not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(AuthenticateUser(true)).
		To(GetMoviesByUser))*/

	return
}

func (ws handler) GenerateToken(req *restful.Request, res *restful.Response) {
	auth := req.HeaderParameter(constant.AuthenticationHeaderName)
	token, err := ws.service.GenerateTokenFromAuthenticationHeader(auth)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, token)
}

func (ws handler) CreateUser(req *restful.Request, res *restful.Response) {
	var userCreation Creation
	err := req.ReadEntity(&userCreation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}
	view, err := ws.service.Create(userCreation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, view)
}

func (ws handler) Update(req *restful.Request, res *restful.Response) {
	var update Update
	err := req.ReadEntity(&update)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}

	user, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	update.UserID = user.ID

	view, err := ws.service.Update(update)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (ws handler) UpdatePassword(req *restful.Request, res *restful.Response) {
	var updatePassword UpdatePassword
	err := req.ReadEntity(&updatePassword)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}

	user, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	updatePassword.UserID = user.ID

	view, err := ws.service.UpdatePassword(updatePassword)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (ws handler) Delete(req *restful.Request, res *restful.Response) {
	var deletion Delete
	user, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	deletion.UserID = user.ID

	err = ws.service.Delete(deletion)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusNoContent, "")
}

func (ws handler) Get(req *restful.Request, res *restful.Response) {
	userID, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	user, err := ws.service.Get(Get{UserID: userID})
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusFound, user)
}

func (ws handler) GetOwnUserInfo(req *restful.Request, res *restful.Response) {
	user, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	view, err := ws.service.GetOwnInfo(user.ID)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (ws handler) UsernameExists(req *restful.Request, res *restful.Response) {
	username := req.PathParameter("username")

	user, err := ws.service.Get(Get{Username: username})

	var exists bool
	if err != nil {
		exists = false
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Sugar.Errorf("Cannot check if username exists, got error : %s", err.Error())
		}
	} else {
		exists = user.Username == username
	}
	if exists {
		res.WriteHeaderAndEntity(http.StatusFound, exists)
		return
	} else {
		res.WriteHeaderAndEntity(http.StatusNotFound, exists)
		return
	}
}

func (ws handler) SearchUsers(req *restful.Request, res *restful.Response) {
	keyword := req.QueryParameter("search")

	filters := Filters{
		Keyword: keyword,
	}

	views, err := ws.service.Search(filters)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, views)
}
