package api

import (
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

func CreateCircle(req *restful.Request, res *restful.Response) {
	var circle model.Circle
	err := req.ReadEntity(&circle)
	if err != nil {
		res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
			model.ErrInternalApiUnprocessableEntity.CodeError())
		return
	}
	err2, newCircle := service.CreateOrUpdateCircle(circle)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, newCircle)
}

func GetCircles(req *restful.Request, res *restful.Response) {
	name := req.QueryParameter("name")
	username := req.HeaderParameter("username")
	err, circles := service.GetCircles(username, "name LIKE ?", "%" + name + "%")
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circles)
}

func DeleteCircle(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	if circleIdStr != "" {
		circleId, err := strconv.Atoi(circleIdStr)
		if err != nil {
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		err2 := service.DeleteCircle(uint(circleId))
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func UpdateCircle(req *restful.Request, res *restful.Response) {
	circleId := req.PathParameter("circleId")
	var circle model.Circle
	if circleId != "" {
		err := req.ReadEntity(&circle)
		if err != nil {
			res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
				model.ErrInternalApiUnprocessableEntity.CodeError())
			return
		}
		id, err := strconv.Atoi(circleId)
		if err != nil {
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		circle.ID = uint(id)
		var err2 model.CustomError
		err2, circle = service.CreateOrUpdateCircle(circle)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, circle)
}

func GetCircle(req *restful.Request, res *restful.Response) {
	circleId := req.PathParameter("circleId")
	var circles []model.Circle
	username := req.HeaderParameter("username")
	if circleId != "" {
		var err2 model.CustomError
		err2, circles = service.GetCircles(username, "id = ?", circleId)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	if len(circles) != 1 {
		res.WriteHeaderAndEntity(model.ErrInternalDatabaseResourceNotFound.HttpCode(), model.ErrInternalDatabaseResourceNotFound.CodeError())
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
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		var err2 model.CustomError
		err2, circle = service.AddUserToCircle(uint(circleId), uint(userId))
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
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
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		var err2 model.CustomError
		err2, circle = service.RemoveUserFromCircle(uint(circleId), uint(userId))
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
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
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		var err2 model.CustomError
		err2, movies = service.GetMoviesForCircle(uint(circleId), sortParameter)
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}

func GetMovieOfCircle(req *restful.Request, res *restful.Response) {
	circleIdStr := req.PathParameter("circleId")
	movieIdStr := req.PathParameter("movieId")
	var movieFound model.Movie
	if circleIdStr != "" && movieIdStr != "" {
		circleId, err := strconv.Atoi(circleIdStr)
		if err != nil {
			model.NewCustomError(err, model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequestCode)
			return
		}
		err2, movies := service.GetMoviesForCircle(uint(circleId), "date:desc")
		if err2.IsNotNil() {
			res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
			return
		}
		for _, movie := range movies {
			if movie.ID == movieIdStr {
				movieFound = movie
			}
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movieFound)
}