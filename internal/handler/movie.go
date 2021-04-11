package handler

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/domain/movieDom"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"net/http"
)

type movieHandler struct {
	service movieDom.Service
}

func NewMovieHandler(svc movieDom.Service) movieHandler {
	return movieHandler{
		service:    svc,
	}
}

func (api movieHandler) WebService() *restful.WebService {
	wsMovie := &restful.WebService{}
	wsMovie.Path("/v1/movies")

	wsMovie.Route(wsMovie.GET("/").
		Doc("Get movie or series by search").
		Param(wsMovie.QueryParameter("title", "Get movie or series by title").DataType("string")).
		Param(wsMovie.QueryParameter("type", "Type of media to search (movie, series, episode)").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "Movie not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(FindMovies))

	wsMovie.Route(wsMovie.GET("/{movieId}").
		Doc("Get movie by ID").
		Param(wsMovie.PathParameter("id", "Get movie by ID (based on IMDb ids)").DataType("string")).
		Writes(model.Movie{}).
		Returns(200, "OK", model.Movie{}).
		Returns(404, "Movie not found", typedErrors.ErrRepositoryResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetMovieById))

	return wsMovie
}

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
	_, username := service.CheckTokenAndGetUsername(req)
	var movie model.Movie
	var err typedErrors.CustomError

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
		res.WriteHeaderAndEntity(typedErrors.ErrApiBadRequest.HttpCode(), typedErrors.ErrApiBadRequest.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movie)
}
