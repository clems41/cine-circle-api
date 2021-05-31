package circleDom

import (
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	webServicePkg "cine-circle/internal/webService"
	"github.com/emicklei/go-restful"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (api handler) WebServices() (handlers []*restful.WebService) {
	wsCircle := &restful.WebService{}
	handlers = append(handlers, wsCircle)

	wsCircle.Path("/v1/circles")

	wsCircle.Route(wsCircle.POST("/").
		Doc("Create new circle").
		Reads(Creation{}).
		Returns(http.StatusCreated, "Created", View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Create))

	wsCircle.Route(wsCircle.PUT("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Update existing circle").
		Reads(Update{}).
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Update))

	wsCircle.Route(wsCircle.GET("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Get existing circle").
		Returns(http.StatusFound, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Get))

	wsCircle.Route(wsCircle.DELETE("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to delete").DataType("int")).
		Doc("Delete existing circle").
		Returns(http.StatusNoContent, "Deleted", nil).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Delete))

	wsCircle.Route(wsCircle.PUT("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "userID to add at the circle").DataType("int")).
		Doc("Add user to circle").
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.AddUser))

	wsCircle.Route(wsCircle.DELETE("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "userID to delete from the circle").DataType("int")).
		Doc("Delete user from circle").
		Returns(http.StatusOK, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.DeleteUser))

	return
}

func (api handler) Create(req *restful.Request, res *restful.Response) {
	var creation Creation
	err := req.ReadEntity(&creation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}

	// Add automatically creator into circle
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	creation.UserIDFromRequest = userFromRequest.ID

	view, err := api.service.Create(creation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, view)
}

func (api handler) Update(req *restful.Request, res *restful.Response) {
	circleID, err := utils.StrToID(req.PathParameter("circleId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	var update Update
	err = req.ReadEntity(&update)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}

	// Check if user sending request is part of the circle
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	update.UserIDFromRequest = userFromRequest.ID
	update.CircleID = circleID

	view, err := api.service.Update(update)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (api handler) Delete(req *restful.Request, res *restful.Response) {
	circleID, err := utils.StrToID(req.PathParameter("circleId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	// Check if user sending request is part of the circle
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	deletion := Deletion{
		CircleID:          circleID,
		UserIDFromRequest: userFromRequest.ID,
	}
	err = api.service.Delete(deletion)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusNoContent, "")
}

func (api handler) Get(req *restful.Request, res *restful.Response) {
	circleID, err := utils.StrToID(req.PathParameter("circleId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	// Check if user sending request is part of the circle
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	get := Get{
		CircleID:          circleID,
		UserIDFromRequest: userFromRequest.ID,
	}
	view, err := api.service.Get(get)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusFound, view)
}

func (api handler) AddUser(req *restful.Request, res *restful.Response) {
	circleID, err := utils.StrToID(req.PathParameter("circleId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	userIDToAdd, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	updateUser := UpdateUser{
		CircleID:          circleID,
		UserIDToUpdate:    userIDToAdd,
		UserIDFromRequest: userFromRequest.ID,
	}
	view, err := api.service.AddUser(updateUser)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (api handler) DeleteUser(req *restful.Request, res *restful.Response) {
	circleID, err := utils.StrToID(req.PathParameter("circleId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	userIDToAdd, err := utils.StrToID(req.PathParameter("userId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	updateUser := UpdateUser{
		CircleID:          circleID,
		UserIDToUpdate:    userIDToAdd,
		UserIDFromRequest: userFromRequest.ID,
	}
	view, err := api.service.DeleteUser(updateUser)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}
