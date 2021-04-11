package handler

import (
	"cine-circle/internal/domain/watchlistDom"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"cine-circle/internal/typedErrors"
	"github.com/emicklei/go-restful"
	"net/http"
)

type watchlistHandler struct {
	service watchlistDom.Service
}

func NewWatchlistHandler(svc watchlistDom.Service) watchlistHandler {
	return watchlistHandler{
		service:    svc,
	}
}

func (api watchlistHandler) WebService() *restful.WebService {
	wsWatchlist := &restful.WebService{}
	wsWatchlist.Path("/v1/watchlist")

	wsWatchlist.Route(wsWatchlist.POST("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to add in watchlist").DataType("int")).
		Doc("Add movie to user's watchlist").
		Writes("").
		Returns(201, "Created", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(AddToWatchlist))

	wsWatchlist.Route(wsWatchlist.DELETE("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to remove from watchlist").DataType("int")).
		Doc("remove movie from users' watchlist").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(RemoveFromWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/").
		Doc("get movies from users' watchlist").
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to check").DataType("int")).
		Doc("Know if movie is already in user's watchlist").
		Writes([]model.Movie{}).
		Returns(200, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(IsInWatchlist))

	return wsWatchlist
}

func AddToWatchlist(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	_, username := service.CheckTokenAndGetUsername(req)
	err := service.AddMovieToWatchlist(username, movieId)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func RemoveFromWatchlist(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	_, username := service.CheckTokenAndGetUsername(req)
	err := service.RemoveMovieFromWatchlist(username, movieId)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func GetWatchlist(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
	err, movies := service.GetMoviesFromWatchlist(username)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movies)
}

func IsInWatchlist(req *restful.Request, res *restful.Response) {
	_, username := service.CheckTokenAndGetUsername(req)
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

