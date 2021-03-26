package api

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func FindMovies(req *restful.Request, res *restful.Response) {
	title := req.QueryParameter("title")
	mediaType := req.QueryParameter("type")
	err, movieSearch := omdb.FindMovieBySearch(title, mediaType)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movieSearch)
}

func GetMovieById(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	username := req.HeaderParameter("username")
	var movie model.Movie
	var err model.CustomError

	if movieId !=  "" {
		err, movie = omdb.FindMovieByID(movieId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
		err = service.AddUserRatings(username, &movie)
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
