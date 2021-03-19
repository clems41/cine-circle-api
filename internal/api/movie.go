package api

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/model"
	"github.com/emicklei/go-restful"
	"net/http"
)

func FindMovie(req *restful.Request, res *restful.Response) {
	movieId := req.QueryParameter("id")
	title := req.QueryParameter("title")
	var movie model.Movie
	var err model.CustomError

	if movieId !=  "" {
		err, movie = omdb.FindMovieByID(movieId)
		if err.IsNotNil() {
			res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
			return
		}
	} else if title != "" {
		err, movie = omdb.FindMovieByTitle(title)
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
