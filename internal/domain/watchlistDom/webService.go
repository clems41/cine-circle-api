package watchlistDom

import (
	"github.com/emicklei/go-restful"
)

type watchlistHandler struct {
	service Service
}

func NewWatchlistHandler(svc Service) *watchlistHandler {
	return &watchlistHandler{
		service:    svc,
	}
}

func (api watchlistHandler) WebService() *restful.WebService {
	wsWatchlist := &restful.WebService{}
	wsWatchlist.Path("/v1/watchlist")

/*	wsWatchlist.Route(wsWatchlist.POST("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to add in watchlist").DataType("int")).
		Doc("Add movie to user's watchlist").
		Returns(201, "Created", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(AddToWatchlist))

	wsWatchlist.Route(wsWatchlist.DELETE("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to remove from watchlist").DataType("int")).
		Doc("remove movie from users' watchlist").
		Returns(http.StatusOK, "OK", "").
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(RemoveFromWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/").
		Doc("get movies from users' watchlist").
		Reads(model.MovieSearch{}).
		Returns(http.StatusOK, "OK", model.MovieSearch{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(GetWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to check").DataType("int")).
		Doc("Know if movie is already in user's watchlist").
		Reads([]model.Movie{}).
		Returns(http.StatusOK, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Filter(authenticateUser(true)).
		To(IsInWatchlist))*/

	return wsWatchlist
}
