package handler

import (
	"cine-circle/internal/domain/circleDom"
	"github.com/emicklei/go-restful"
)

type circleHandler struct {
	service circleDom.Service
}

func NewCircleHandler(svc circleDom.Service) *circleHandler {
	return &circleHandler{
		service:    svc,
	}
}

func (api circleHandler) WebService() *restful.WebService {
	wsCircle := &restful.WebService{}
	wsCircle.Path("/v1/circles")

/*	wsCircle.Route(wsCircle.POST("/").
		Doc("Create new circle").
		Writes(model.Circle{}).
		Returns(201, "Created", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(CreateCircle))

	wsCircle.Route(wsCircle.PUT("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Update existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(UpdateCircle))

	wsCircle.Route(wsCircle.GET("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Get existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(GetCircle))

	wsCircle.Route(wsCircle.GET("/").
		Param(wsCircle.QueryParameter("name", "find circles by name").DataType("string")).
		Doc("Search for circles").
		Writes([]model.Circle{}).
		Returns(200, "Found", []model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(GetCircles))

	wsCircle.Route(wsCircle.GET("/{circleId}/movies").
		Param(wsCircle.PathParameter("circleId", "ID of circle to get movies").DataType("int")).
		Param(wsCircle.QueryParameter("sort", "way of sorting movies").DataType("string")).
		Doc("Get movies of circle with sorting (default='date:desc'").
		Writes([]model.Movie{}).
		Returns(200, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(GetMoviesOfCircle))

	wsCircle.Route(wsCircle.DELETE("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to delete").DataType("int")).
		Doc("Delete existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(DeleteCircle))

	wsCircle.Route(wsCircle.PUT("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "ID of user to add to circle").DataType("int")).
		Doc("Add user to existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(AddUserToCircle))

	wsCircle.Route(wsCircle.DELETE("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "ID of user to remove from circle").DataType("int")).
		Doc("Remove user from existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(RemoveUserFromCircle))*/

	return wsCircle
}