package recommendationDom

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

func (api handler) WebServices() (webServices []*restful.WebService) {
	wsReco := &restful.WebService{}
	webServices = append(webServices, wsReco)

	wsReco.Path("/v1/recommendations")

	wsReco.Route(wsReco.POST("/").
		Doc("Send new recommendation").
		Reads(Creation{}).
		Returns(http.StatusCreated, "Created", nil).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, "Not processable, impossible to serialize json", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Create))

	wsReco.Route(wsReco.GET("/").
		Param(wsReco.QueryParameter("page", "num of page to get").DataType("int").DefaultValue("1")).
		Param(wsReco.QueryParameter("pageSize", "number of element if one page").DataType("int").DefaultValue("10")).
		Param(wsReco.QueryParameter("sort", "way of sorting elements (date:asc)").DataType("string").DefaultValue("date:desc")).
		Param(wsReco.QueryParameter("recommendationType", "filter on type (received, sent or both)").DataType("string").DefaultValue("received")).
		Param(wsReco.QueryParameter("movieId", "get only recommendations for specific movie").DataType("int").DefaultValue("")).
		Param(wsReco.QueryParameter("circleId", "get only recommendations for specific circle").DataType("int").DefaultValue("")).
		Doc("List, filter and sort recommendations").
		Returns(http.StatusOK, "Created", ViewList{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.List))

	wsReco.Route(wsReco.GET("/users/").
		Param(wsReco.QueryParameter("page", "num of page to get").DataType("int")).
		Param(wsReco.QueryParameter("pageSize", "number of element if one page").DataType("int")).
		Doc("List all users that can received recommendation from user").
		Returns(http.StatusOK, "Created", ViewList{}).
		Returns(http.StatusBadRequest, "Bad request, fields not validated", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.ListUsers))

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
	userFromRequest, err := webServicePkg.WhoAmI(req)
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

func (api handler) List(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	// get query parameters for filtering search
	var filters Filters
	filters.PaginationRequest, err = utils.ExtractPaginationRequest(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	filters.SortingRequest, err = utils.ExtractSortingRequest(req, "date", true)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	filters.UserID = userFromRequest.ID
	filters.RecommendationType = req.QueryParameter("recommendationType")

	movieIdStr := req.QueryParameter("movieId")
	if movieIdStr != "" {
		filters.MovieID, err =  utils.StrToID(movieIdStr)
		if err != nil {
			webServicePkg.HandleHTTPError(req, res, err)
			return
		}
	}

	circleIdStr := req.QueryParameter("circleId")
	if circleIdStr != "" {
		filters.CircleID, err =  utils.StrToID(circleIdStr)
		if err != nil {
			webServicePkg.HandleHTTPError(req, res, err)
			return
		}
	}

	view, err := api.service.List(filters)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (api handler) ListUsers(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	var usersFilters UsersFilters
	usersFilters.PaginationRequest, err = utils.ExtractPaginationRequest(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	usersFilters.UserID = userFromRequest.ID

	list, err := api.service.ListUsers(usersFilters)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, list)
}
