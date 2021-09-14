package movieDom

import (
	utils2 "cine-circle/pkg/utils"
	"cine-circle/pkg/webService"
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
	wsMovie := &restful.WebService{}
	webServices = append(webServices, wsMovie)

	wsMovie.Path("/v1/movies")

	wsMovie.Route(wsMovie.GET("/{movieId}").
		Param(wsMovie.PathParameter("movieId", "ID of movie").DataType("int")).
		Doc("Get movie").
		Returns(http.StatusFound, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webService.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webService.FormattedJsonError{}).
		Filter(webService.LogRequest()).
		Filter(webService.AuthenticateUser()).
		To(api.Get))

	wsMovie.Route(wsMovie.GET("/").
		Param(wsMovie.QueryParameter("page", "num of page to get").DataType("int")).
		Param(wsMovie.QueryParameter("query", "query to search among tv shows and movies").DataType("string")).
		Doc("Search movies").
		Returns(http.StatusOK, "OK", SearchView{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webService.FormattedJsonError{}).
		Filter(webService.LogRequest()).
		Filter(webService.AuthenticateUser()).
		To(api.Search))

	return
}

func (api handler) Get(req *restful.Request, res *restful.Response) {
	movieID, err := utils2.StrToID(req.PathParameter("movieId"))
	if err != nil {
		webService.HandleHTTPError(req, res, err)
		return
	}

	view, err := api.service.Get(movieID)
	if err != nil {
		webService.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusFound, view)
}

func (api handler) Search(req *restful.Request, res *restful.Response) {
	var filters Filters
	var err error
	filters.PaginationRequest, err = utils2.ExtractPaginationRequest(req)
	if err != nil {
		webService.HandleHTTPError(req, res, err)
		return
	}

	filters.Query = req.QueryParameter("query")

	result, err := api.service.Search(filters)
	if err != nil {
		webService.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, result)
}
