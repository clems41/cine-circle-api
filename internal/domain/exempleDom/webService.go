package exempleDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/constant/webserviceConst"
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/httpServer/httpError"
	"cine-circle-api/pkg/httpServer/middleware"
	"cine-circle-api/pkg/utils/httpUtils"
	"cine-circle-api/pkg/utils/idUtils"
	"cine-circle-api/pkg/utils/validationUtils"
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"net/http"
	"strings"
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
	wsUser := new(restful.WebService)
	tags := []string{swaggerConst.UserTag}

	wsUser.Path(basePath)

	wsUser.Route(wsUser.POST("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Créer un nouvel exemple").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(CreateForm{}).
		Writes(CreateView{}).
		Returns(http.StatusCreated, "Exemple créé", CreateView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Create))

	wsUser.Route(wsUser.PUT("/{exempleId}").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Modification des informations d'un exemple par un administrateur").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(wsUser.PathParameter("exempleId", "ID de l'exemple à modifier").
			DataType("int").Required(true)).
		Reads(UpdateForm{}).
		Writes(UpdateView{}).
		Returns(http.StatusOK, "Exemple mis à jour", UpdateView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Update))

	wsUser.Route(wsUser.DELETE("/{exempleId}").
		Doc("Suppression d'un exemple").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(wsUser.PathParameter("exempleId", "ID de l'exemple à supprimer").
			DataType("int").Required(true)).
		Returns(http.StatusOK, "Exemple supprimé", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Delete))

	wsUser.Route(wsUser.GET("/{exempleId}").
		Doc("Récupérer les informations d'un exemple").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(wsUser.PathParameter("exempleId", "ID de l'exemple à récupérer").
			DataType("int").Required(true)).
		Returns(http.StatusOK, "Informations récupérées", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Get))

	wsUser.Route(wsUser.GET("/").
		Produces(restful.MIME_JSON).
		Param(wsUser.QueryParameter(webserviceConst.PageQueryName, webserviceConst.PageQueryDescription).
			DataType(webserviceConst.PageQueryType).
			DefaultValue(fmt.Sprintf("%d", defaultPage))).
		Param(wsUser.QueryParameter(webserviceConst.PageSizeQueryName, webserviceConst.PageSizeQueryDescription).
			DataType(webserviceConst.PageSizeQueryType).
			DefaultValue(fmt.Sprintf("%d", defaultPageSize))).
		Param(wsUser.QueryParameter(webserviceConst.SortQueryName, webserviceConst.SortQueryDescription).
			DataType(webserviceConst.SortQueryType).
			DefaultValue(strings.Join(defaultSort, ","))).
		Doc("Retourne la liste des exemples de l'application avec pagination et tri par champ").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(SearchView{}).
		Returns(http.StatusOK, "Liste récupérée", SearchView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Search))

	return wsUser
}

func (hd *handler) Create(req *restful.Request, res *restful.Response) {
	var form CreateForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	view, err := hd.service.Create(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, view)
}

func (hd *handler) Update(req *restful.Request, res *restful.Response) {
	var form UpdateForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	exempleId, err := idUtils.StrToID(req.PathParameter("exempleId"))
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form.ExempleId = exempleId

	view, err := hd.service.Update(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Delete(req *restful.Request, res *restful.Response) {
	exempleId, err := idUtils.StrToID(req.PathParameter("exempleId"))
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form := DeleteForm{ExempleId: exempleId}

	err = hd.service.Delete(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) Get(req *restful.Request, res *restful.Response) {
	exempleId, err := idUtils.StrToID(req.PathParameter("exempleId"))
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form := GetForm{ExempleId: exempleId}

	view, err := hd.service.Get(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Search(req *restful.Request, res *restful.Response) {
	var form SearchForm
	err := httpUtils.UnmarshallQueryParameters(req, &form, defaultSearchQueryParametersValues)
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
