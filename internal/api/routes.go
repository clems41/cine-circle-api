package api

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func DefineRoutes() []*restful.WebService {
	wsRoot := &restful.WebService{}

	// HEALTH
	wsRoot.Route(wsRoot.GET("/health/ok").
		Doc("Simple API health check").
		To(func(req *restful.Request, res *restful.Response) {
			res.WriteHeader(http.StatusOK)
		}))

	// MOVIE

	wsMovie := &restful.WebService{}
	wsMovie.Path("/v1/movies")

	wsMovie.Route(wsMovie.GET("/").
		Doc("Get movie or series by search").
		Param(wsMovie.QueryParameter("title", "Get movie or series by title").DataType("string")).
		Param(wsMovie.QueryParameter("type", "Type of media to search (movie, series, episode)").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "Movie not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(FindMovies))

	wsMovie.Route(wsMovie.GET("/{movieId}").
		Doc("Get movie by ID").
		Param(wsMovie.PathParameter("id", "Get movie by ID (based on IMDb ids)").DataType("string")).
		Writes(model.Movie{}).
		Returns(200, "OK", model.Movie{}).
		Returns(404, "Movie not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetMovieById))

	// LOGIN

	wsLogin := &restful.WebService{}
	wsLogin.Path("/v1")

	wsLogin.Route(wsLogin.POST("/signup").
		Doc("Create new user").
		Writes(model.User{}).
		Returns(201, "Created", model.User{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(false)).
		To(CreateUser))

	wsLogin.Route(wsLogin.POST("/signin").
		Doc("Connect with existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(false)).
		To(GetToken))


	// USER

	wsUser := &restful.WebService{}
	wsUser.Path("/v1/users")

	wsUser.Route(wsUser.PUT("/{userId}").
		Doc("Update existing user").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(UpdateUser))

	wsUser.Route(wsUser.DELETE("/{userId}").
		Doc("Delete existing user").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(DeleteUser))

	wsUser.Route(wsUser.GET("/").
		Doc("Search for user(s)").
		Param(wsUser.QueryParameter("username", "search user by username").DataType("string")).
		Param(wsUser.QueryParameter("email", "search user by email").DataType("string")).
		Param(wsUser.QueryParameter("fullname", "search user by fullname").DataType("string")).
		Writes([]model.User{}).
		Returns(200, "OK", []model.User{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(SearchUsers))

	wsUser.Route(wsUser.GET("/{userId}").
		Param(wsUser.PathParameter("userId", "ID of sought user").DataType("int")).
		Doc("Get user info from ID").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetUser))

	wsUser.Route(wsUser.GET("/me").
		Doc("Get user info from token").
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetOwnUserInfo))

	wsUser.Route(wsUser.GET("/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUser.PathParameter("username", "username of sought user").DataType("string")).
		Writes("").
		Returns(200, "OK", "").
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(false)).
		To(UsernameExists))

	wsUser.Route(wsUser.GET("/{userId}/movies").
		Doc("Get all movies that user had rated").
		Param(wsUser.PathParameter("userId", "username of sought user").DataType("string")).
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(true)).
		To(GetMoviesByUser))

	// RATING

	wsRating := &restful.WebService{}
	wsRating.Path("/v1/ratings")

	wsRating.Route(wsRating.POST("/{movieId}").
		Param(wsRating.PathParameter("movieId", "ID of the movie to rate").DataType("int")).
		Doc("Add rating to movie for specific user").
		Writes(model.Rating{}).
		Returns(201, "Created", model.Rating{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Rating",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(AddRating))

	// WATCHLIST

	wsWatchlist := &restful.WebService{}
	wsWatchlist.Path("/v1/watchlist")

	wsWatchlist.Route(wsWatchlist.POST("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to add in watchlist").DataType("int")).
		Doc("Add movie to user's watchlist").
		Writes("").
		Returns(201, "Created", "").
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(AddToWatchlist))

	wsWatchlist.Route(wsWatchlist.DELETE("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to remove from watchlist").DataType("int")).
		Doc("remove movie from users' watchlist").
		Writes("").
		Returns(200, "OK", "").
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(RemoveFromWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/").
		Doc("get movies from users' watchlist").
		Writes(model.MovieSearch{}).
		Returns(200, "OK", model.MovieSearch{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetWatchlist))

	wsWatchlist.Route(wsWatchlist.GET("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to check").DataType("int")).
		Doc("Know if movie is already in user's watchlist").
		Writes([]model.Movie{}).
		Returns(200, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(IsInWatchlist))

	// CIRCLE

	wsCircle := &restful.WebService{}
	wsCircle.Path("/v1/circles")

	wsCircle.Route(wsCircle.POST("/").
		Doc("Create new circle").
		Writes(model.Circle{}).
		Returns(201, "Created", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(CreateCircle))

	wsCircle.Route(wsCircle.PUT("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Update existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(UpdateCircle))

	wsCircle.Route(wsCircle.GET("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Doc("Get existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetCircle))

	wsCircle.Route(wsCircle.GET("/").
		Param(wsCircle.QueryParameter("name", "find circles by name").DataType("string")).
		Doc("Search for circles").
		Writes([]model.Circle{}).
		Returns(200, "Found", []model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(GetCircles))

	wsCircle.Route(wsCircle.GET("/{circleId}/movies").
		Param(wsCircle.PathParameter("circleId", "ID of circle to get movies").DataType("int")).
		Param(wsCircle.QueryParameter("sort", "way of sorting movies").DataType("string")).
		Doc("Get movies of circle with sorting (default='date:desc'").
		Writes([]model.Movie{}).
		Returns(200, "OK", []model.Movie{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Filter(filterUser(true)).
		To(GetMoviesOfCircle))

	wsCircle.Route(wsCircle.DELETE("/{circleId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to delete").DataType("int")).
		Doc("Delete existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(DeleteCircle))

	wsCircle.Route(wsCircle.PUT("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "ID of user to add to circle").DataType("int")).
		Doc("Add user to existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(AddUserToCircle))

	wsCircle.Route(wsCircle.DELETE("/{circleId}/{userId}").
		Param(wsCircle.PathParameter("circleId", "ID of circle to update").DataType("int")).
		Param(wsCircle.PathParameter("userId", "ID of user to remove from circle").DataType("int")).
		Doc("Remove user from existing circle").
		Writes(model.Circle{}).
		Returns(200, "Updated", model.Circle{}).
		Returns(400, "Bad request, fields not validated", model.ErrInternalApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Circle",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(RemoveUserFromCircle))

	return []*restful.WebService{wsRoot, wsMovie, wsUser, wsRating, wsCircle, wsWatchlist, wsLogin}
}

// Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func filterUser(needAuthentication bool) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		if needAuthentication {
			logger.Sugar.Debugf("Token will be checked")
			if err, _ := service.CheckTokenAndGetUsername(req); err.IsNotNil() {
				res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
				return
			}
		}
		chain.ProcessFilter(req, res)
	}
	return filter
}