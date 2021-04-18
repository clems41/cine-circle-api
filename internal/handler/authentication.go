package handler

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain/authenticationDom"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"net/http"
)

type authenticationHandler struct {
	service authenticationDom.Service
}

func NewAuthenticationHandler(svc authenticationDom.Service) *authenticationHandler {
	return &authenticationHandler{
		service:    svc,
	}
}

func (handler authenticationHandler) WebService() *restful.WebService {
	wsAuthentication := &restful.WebService{}
	wsAuthentication.Path("/v1")

	wsAuthentication.Route(wsAuthentication.POST("/signup").
		Doc("Create new user (signup)").
		Writes(userDom.Creation{}).
		Returns(201, "Created", userDom.Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json",typedErrors.CustomError{}).
		Filter(logRequest()).
		To(handler.CreateUser))

	wsAuthentication.Route(wsAuthentication.POST("/signin").
		Doc("Generate token from username and password (basic authentication)").
		Writes(nil).
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json", typedErrors.CustomError{}).
		Filter(logRequest()).
		To(handler.GenerateToken))

	return wsAuthentication
}

func (handler authenticationHandler) GenerateToken(req *restful.Request, res *restful.Response) {
	auth := req.HeaderParameter(constant.AuthenticationHeaderName)
	token, err := handler.service.GenerateTokenFromAuthenticationHeader(auth)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, token)
}

func (handler authenticationHandler) CreateUser(req *restful.Request, res *restful.Response) {
	var userCreation userDom.Creation
	err := req.ReadEntity(&userCreation)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	user, err := handler.service.CreateUser(userCreation)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, user)
}
