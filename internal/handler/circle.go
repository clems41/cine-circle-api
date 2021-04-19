package handler

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/domain/circleDom"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

type circleHandler struct {
	service circleDom.Service
}

func NewCircleHandler(svc circleDom.Service) *circleHandler {
	return &circleHandler{
		service:    svc,
	}
}

func (api circleHandler) WebService() *restful.WebService {
	wsCircle := &restful.WebService{}
	wsCircle.Path("/v1/circles")

	wsCircle.Route(wsCircle.POST("/").
		Doc("Create new circle").
		Writes(circleDom.Creation{}).
		Returns(201, "Created", circleDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Create))

	wsCircle.Route(wsCircle.PUT("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Update existing circle").
		Writes(circleDom.Update{}).
		Returns(200, "OK", circleDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Update))

	wsCircle.Route(wsCircle.GET("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Get existing circle").
		Writes(nil).
		Returns(200, "OK", circleDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Get))

	wsCircle.Route(wsCircle.DELETE("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to delete").DataType("int")).
		Doc("Delete existing circle").
		Writes(nil).
		Returns(200, "OK", circleDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.Delete))

	return wsCircle
}

func (api circleHandler) Create(req *restful.Request, res *restful.Response) {
	var creation circleDom.Creation
	err := req.ReadEntity(&creation)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	// get current user from token
	user, err := CommonHandler.WhoAmI(req)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	// add current user if not already in circle's users
	if !utils.ContainsID(creation.UsersID, user.UserID) {
		creation.UsersID = append(creation.UsersID, user.UserID)
	}

	circle, err := api.service.Create(creation)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func (api circleHandler) Update(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	circleId, err := strconv.Atoi(circleIdStr)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	var update circleDom.Update
	err = req.ReadEntity(&update)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	// get current user from token
	user, err := CommonHandler.WhoAmI(req)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	// add current user if not already in circle's users
	if !utils.ContainsID(update.UsersID, user.UserID) {
		update.UsersID = append(update.UsersID, user.UserID)
	}

	update.CircleID = domain.IDType(circleId)

	circle, err := api.service.Update(update)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func (api circleHandler) Delete(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	circleId, err := strconv.Atoi(circleIdStr)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	deleteCircle := circleDom.Delete{CircleID: domain.IDType(circleId)}
	err = api.service.Delete(deleteCircle)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (api circleHandler) Get(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	circleId, err := strconv.Atoi(circleIdStr)
	if err != nil {
		handleHTTPError(res, typedErrors.NewApiBadRequestErrorf(err.Error()))
		return
	}

	get := circleDom.Get{CircleID: domain.IDType(circleId)}
	circle, err := api.service.Get(get)
	if err != nil {
		handleHTTPError(res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, circle)
}