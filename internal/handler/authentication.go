package handler

import (
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"net/http"
)

func NewAuthenticationHandler() *restful.WebService {
	wsAuthentication := &restful.WebService{}
	wsAuthentication.Path("/v1")

	wsAuthentication.Route(wsAuthentication.POST("/signup").
		Doc("Create new user").
		Writes(model.User{}).
		Returns(201, "Created", model.User{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(false)).
		To(CreateUser))

	wsAuthentication.Route(wsAuthentication.POST("/signin").
		Doc("Connect with existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(false)).
		To(GetToken))

	return wsAuthentication
}

func GetToken(req *restful.Request, res *restful.Response) {
	auth := req.HeaderParameter("Authorization")
	err, token := service.GetTokenFromAuthentication(auth)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, token)
}
