package mediaDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/constant/webserviceConst"
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/httpServer/httpError"
	"cine-circle-api/pkg/httpServer/middleware"
	"cine-circle-api/pkg/utils/httpUtils"
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
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

func (hd *handler) WebService() *restful.WebService {
	wsMedia := new(restful.WebService)
	tags := []string{swaggerConst.MediaTag}

	wsMedia.Path(basePath)

	wsMedia.Route(wsMedia.GET(fmt.Sprintf("/%s", mediaIdParameter.Joker())).
		Doc("Get media (movie or tv show) by id").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(mediaIdParameter.PathParameter()).
		Param(languageParameter.QueryParameter()).
		Returns(http.StatusOK, "Media found", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusForbidden, webserviceConst.ForbiddenMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Get))

	wsMedia.Route(wsMedia.GET("/").
		Produces(restful.MIME_JSON).
		Param(languageParameter.QueryParameter()).
		Param(keywordParameter.QueryParameter()).
		Doc("Search among medias ones that are matching with keyword").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(SearchView{}).
		Returns(http.StatusOK, "Search OK", SearchView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusForbidden, webserviceConst.ForbiddenMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Search))

	return wsMedia
}

func (hd *handler) Get(req *restful.Request, res *restful.Response) {
	mediaId, err := mediaIdParameter.GetValueFromPathParameter(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	var form GetForm
	err = httpUtils.UnmarshallQueryParameters(req, &form, defaultQueryParametersValues)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form.MediaId = mediaId

	view, err := hd.service.Get(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Search(req *restful.Request, res *restful.Response) {
	var form SearchForm
	err := httpUtils.UnmarshallQueryParameters(req, &form, defaultQueryParametersValues)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	view, err := hd.service.Search(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}
