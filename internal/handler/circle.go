package handler

import (
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

func NewCircleHandler() *restful.WebService {
	wsCircle := &restful.WebService{}
	wsCircle.Path("/v1/circles")

	wsCircle.Route(wsCircle.POST("/").
		Doc("Create new circle").
		Writes(model.Circle{}).
		Returns(201, "Created", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(CreateCircle))

	wsCircle.Route(wsCircle.PUT("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Update existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(UpdateCircle))

	wsCircle.Route(wsCircle.GET("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Get existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetCircle))

	wsCircle.Route(wsCircle.GET("/").
		Param(wsCircle.QueryParameter("name", "find circles by name").DataType("string")).
		Doc("Search for circles").
		Writes([]model.Circle{}).
		Returns(200, "Found", []model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(GetCircles))

	wsCircle.Route(wsCircle.GET("/{circleId}/movies").
		Param(wsCircle.PathParameter("circleId", "ID of circle to get movies").DataType("int")).
		Param(wsCircle.QueryParameter("sort", "way of sorting movies").DataType("string")).
		Doc("Get movies of circle with sorting (default='date:desc'").
		Writes([]model.Movie{}).
		Returns(200, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetMoviesOfCircle))

	wsCircle.Route(wsCircle.DELETE("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to delete").DataType("int")).
		Doc("Delete existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
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
		Filter(filterUser(true)).
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
		Filter(filterUser(true)).
		To(RemoveUserFromCircle))

	return wsCircle
}

func CreateCircle(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
	var circle model.Circle
	err := req.ReadEntity(&circle)
	if err != nil {
		res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
			typedErrors.ErrApiUnprocessableEntity.CodeError())
		return
	}
	err2, newCircle := service.CreateCircle(circle, username)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, newCircle)
}

func GetCircles(req *restful.Request, res *restful.Response) {
	name := req.QueryParameter("name")
	_, username := service.CheckTokenAndGetUsername(req)
	err, circles := service.GetCircles(username, "name LIKE ?", "%" + name + "%")
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circles)
}

func DeleteCircle(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
	circleId := req.PathParameter("circleId")
	if circleId != "" {
		err2 := service.DeleteCircle(circleId, username)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func UpdateCircle(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
	circleId := req.PathParameter("circleId")
	var circle model.Circle
	if circleId != "" {
		err := req.ReadEntity(&circle)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
				typedErrors.ErrApiUnprocessableEntity.CodeError())
			return
		}
		var err2 typedErrors.CustomError
		err2, circle = service.UpdateCircle(circle, circleId, username)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func GetCircle(req *restful.Request, res *restful.Response) {
	circleId := req.PathParameter("circleId")
	var circles []model.Circle
	_, username := service.CheckTokenAndGetUsername(req)
	if circleId != "" {
		var err2 typedErrors.CustomError
		err2, circles = service.GetCircles(username, "id = ?", circleId)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	if len(circles) != 1 {
		res.WriteHeaderAndEntity(typedErrors.ErrRepositoryResourceNotFound.HttpCode(), typedErrors.ErrRepositoryResourceNotFound.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circles[0])
}

func AddUserToCircle(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	userIdStr := req.PathParameter("userId")
	var circle model.Circle
	if circleIdStr != "" && userIdStr != "" {
		circleId, err := strconv.Atoi(circleIdStr)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.NewApiBadRequestError(err).HttpCode(), typedErrors.NewApiBadRequestError(err).CodeError())
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.NewApiBadRequestError(err).HttpCode(), typedErrors.NewApiBadRequestError(err).CodeError())
			return
		}
		var err2 typedErrors.CustomError
		err2, circle = service.AddUserToCircle(uint(circleId), uint(userId))
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func RemoveUserFromCircle(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	userIdStr := req.PathParameter("userId")
	var circle model.Circle
	if circleIdStr != "" && userIdStr != "" {
		circleId, err := strconv.Atoi(circleIdStr)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.NewApiBadRequestError(err).HttpCode(), typedErrors.NewApiBadRequestError(err).CodeError())
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.NewApiBadRequestError(err).HttpCode(), typedErrors.NewApiBadRequestError(err).CodeError())
			return
		}
		var err2 typedErrors.CustomError
		err2, circle = service.RemoveUserFromCircle(uint(circleId), uint(userId))
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func GetMoviesOfCircle(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	sortParameter := req.QueryParameter("sort")
	var movies []model.Movie
	if circleIdStr != "" {
		circleId, err := strconv.Atoi(circleIdStr)
		if err != nil {
			res.WriteHeaderAndEntity(typedErrors.NewApiBadRequestError(err).HttpCode(), typedErrors.NewApiBadRequestError(err).CodeError())
			return
		}
		var err2 typedErrors.CustomError
		err2, movies = service.GetMoviesForCircle(uint(circleId), sortParameter)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}