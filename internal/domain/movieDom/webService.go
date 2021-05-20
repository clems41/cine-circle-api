package movieDom

import (
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
		Writes(nil).
		Returns(http.StatusFound, "OK", View{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.Get))

	return
}

func (api handler) Get(req *restful.Request, res *restful.Response) {
	movieID, err := utils.StrToID(req.PathParameter("movieId"))
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	view, err := api.service.Get(movieID)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusFound, view)
}
