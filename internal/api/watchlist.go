package api

import (
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func AddToWatchlist(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	username := req.HeaderParameter(userHeaderParameter)
	err := service.AddMovieToWatchlist(username, movieId)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func RemoveFromWatchlist(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	username := req.HeaderParameter(userHeaderParameter)
	err := service.RemoveMovieFromWatchlist(username, movieId)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func GetWatchlist(req *restful.Request, res *restful.Response) {
	username := req.HeaderParameter(userHeaderParameter)
	err, movies := service.GetMoviesFromWatchlist(username)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}

func IsInWatchlist(req *restful.Request, res *restful.Response) {
	username := req.HeaderParameter("username")
	movieId := req.PathParameter("movieId")
	err, isIn := service.IsInWatchlist(username, movieId)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	if isIn {
		res.WriteHeaderAndEntity(http.StatusOK, "true")
	} else {
		res.WriteHeaderAndEntity(http.StatusOK, "false")
	}
}

