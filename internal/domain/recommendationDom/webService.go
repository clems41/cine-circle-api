package recommendationDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/constant/webserviceConst"
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/httpServer/authentication"
	"cine-circle-api/pkg/httpServer/httpError"
	"cine-circle-api/pkg/httpServer/middleware"
	"cine-circle-api/pkg/utils/httpUtils"
	"cine-circle-api/pkg/utils/validationUtils"
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
	wsUser := new(restful.WebService)
	tags := []string{swaggerConst.UserTag}

	wsUser.Path(basePath)

	wsUser.Route(wsUser.POST("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Send new recommendation").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(SendForm{}).
		Writes(SearchView{}).
		Returns(http.StatusCreated, "Recommendation sent", SendView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusForbidden, webserviceConst.ForbiddenMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Send))

	wsUser.Route(wsUser.GET("/").
		Produces(restful.MIME_JSON).
		Param(pageQueryParameter.QueryParameter()).
		Param(pageSizeQueryParameter.QueryParameter()).
		Param(mediaIdQueryParameter.QueryParameter()).
		Param(recommendationTypeQueryParameter.QueryParameter()).
		Doc(fmt.Sprintf("List recommendations with pagination, sort and filters")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(SearchView{}).
		Returns(http.StatusOK, "OK", SearchView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusForbidden, webserviceConst.ForbiddenMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Search))

	return wsUser
}

func (hd *handler) Send(req *restful.Request, res *restful.Response) {
	var form SendForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewForbidden().WrapError(err))
		return
	}
	form.SenderId = user.Id

	view, err := hd.service.Send(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, view)
}

func (hd *handler) Search(req *restful.Request, res *restful.Response) {
	var form SearchForm
	err := httpUtils.UnmarshallQueryParameters(req, &form, defaultSearchQueryParameterValues)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewForbidden().WrapError(err))
		return
	}
	form.UserId = user.Id

	view, err := hd.service.Search(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}
