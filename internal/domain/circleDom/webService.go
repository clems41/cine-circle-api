package circleDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/constant/webserviceConst"
	"cine-circle-api/pkg/customError"
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
		Doc("Create new circle").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(CreateForm{}).
		Writes(CreateView{}).
		Returns(http.StatusCreated, "Circle created", CreateView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Create))

	wsUser.Route(wsUser.PUT(fmt.Sprintf("/%s/%s", circleIdPathParameter.Joker(), userIdPathParameter.Joker())).
		Produces(restful.MIME_JSON).
		Doc("Add user into circle").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(circleIdPathParameter.PathParameter()).
		Param(userIdPathParameter.PathParameter()).
		Reads(UpdateForm{}).
		Writes(UpdateView{}).
		Returns(http.StatusOK, "User added", UpdateView{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.AddUser))

	wsUser.Route(wsUser.DELETE(fmt.Sprintf("/%s/%s", circleIdPathParameter.Joker(), userIdPathParameter.Joker())).
		Produces(restful.MIME_JSON).
		Doc("Delete user from circle").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(circleIdPathParameter.PathParameter()).
		Param(userIdPathParameter.PathParameter()).
		Reads(UpdateForm{}).
		Writes(UpdateView{}).
		Returns(http.StatusOK, "User deleted", UpdateView{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.DeleteUser))

	wsUser.Route(wsUser.PUT("/"+circleIdPathParameter.Joker()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Update circle info").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(circleIdPathParameter.PathParameter()).
		Reads(UpdateForm{}).
		Writes(UpdateView{}).
		Returns(http.StatusOK, "Circle updated", UpdateView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Update))

	wsUser.Route(wsUser.DELETE("/"+circleIdPathParameter.Joker()).
		Doc("Delete circle").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(circleIdPathParameter.PathParameter()).
		Returns(http.StatusOK, "Circle deleted", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Delete))

	wsUser.Route(wsUser.GET("/"+circleIdPathParameter.Joker()).
		Doc("Get circle info").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(circleIdPathParameter.PathParameter()).
		Returns(http.StatusOK, "OK", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusNotFound, webserviceConst.NotFoundMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Get))

	wsUser.Route(wsUser.GET("/").
		Produces(restful.MIME_JSON).
		Param(pageQueryParameter.QueryParameter()).
		Param(pageSizeQueryParameter.QueryParameter()).
		Param(circleNameQueryParameter.QueryParameter()).
		Doc(fmt.Sprintf("Search among circles by %s", circleNameQueryParameter.Name)).
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

func (hd *handler) AddUser(req *restful.Request, res *restful.Response) {
	circleId, err := circleIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	userId, err := userIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form := AddUserForm{
		UserId:   userId,
		CircleId: circleId,
	}

	view, err := hd.service.AddUser(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) DeleteUser(req *restful.Request, res *restful.Response) {
	circleId, err := circleIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	userId, err := userIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form := DeleteUserForm{
		UserId:   userId,
		CircleId: circleId,
	}

	view, err := hd.service.DeleteUser(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Update(req *restful.Request, res *restful.Response) {
	var form UpdateForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	circleId, err := circleIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form.CircleId = circleId

	view, err := hd.service.Update(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Delete(req *restful.Request, res *restful.Response) {
	circleId, err := circleIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form := DeleteForm{CircleId: circleId}

	err = hd.service.Delete(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) Get(req *restful.Request, res *restful.Response) {
	circleId, err := circleIdPathParameter.GetValueFromPathParameterAsUint(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form := GetForm{CircleId: circleId}

	view, err := hd.service.Get(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Search(req *restful.Request, res *restful.Response) {
	var form SearchForm
	err := httpUtils.UnmarshallQueryParameters(req, &form, defaultSearchQueryParameterValues)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	view, err := hd.service.Search(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}
