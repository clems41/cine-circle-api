package recommendationDom

import (
	"cine-circle/internal/typedErrors"
	webServicePkg "cine-circle/internal/webService"
	"github.com/emicklei/go-restful"
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

func (api handler) WebServices() (webServices []*restful.WebService) {
	wsReco := &restful.WebService{}
	webServices = append(webServices, wsReco)

	wsReco.Path("/v1/recommendations")

	wsReco.Route(wsReco.POST("/").
		Doc("Send new recommendation").
		Writes(Creation{}).
		Returns(http.StatusCreated, "Created", nil).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Create))

	return
}

func (api handler) Create(req *restful.Request, res *restful.Response) {
	var creation Creation
	err := req.ReadEntity(&creation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, typedErrors.NewUnprocessableEntityErrorf(err.Error()))
		return
	}

	// Get user info from token
	userFromRequest, err := webServicePkg.ActualUserHandler.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	creation.SenderID = userFromRequest.ID

	err = api.service.Create(creation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusCreated, "")
}
