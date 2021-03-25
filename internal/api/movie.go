package api

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func FindMovie(req *restful.Request, res *restful.Response) {
	title := req.QueryParameter("title")
	username := req.HeaderParameter("username")
	var movie model.Movie
	var err model.CustomError

	if title != "" {
		err, movie = omdb.FindMovieByTitle(title)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
		err = service.AddUserRating(username, &movie)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movie)
}

func GetMovieById(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	var movie model.Movie
	var err model.CustomError

	if movieId !=  "" {
		err, movie = omdb.FindMovieByID(movieId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else {
		res.WriteHeaderAndEntity(model.ErrInternalApiBadRequest.HttpCode(), model.ErrInternalApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movie)
}
