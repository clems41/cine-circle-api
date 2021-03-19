package api

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

const (
	userHeaderParameter = "username"
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
		Doc("Get movie by ID or by title").
		Param(wsMovie.QueryParameter("id", "Get movie by ID (based on IMDb ids)").DataType("string")).
		Param(wsMovie.QueryParameter("title", "Get movie by title").DataType("string")).
		Writes(model.Movie{}).
		Returns(200, "OK", model.Movie{}).
		Returns(404, "Movie not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(false)).
		To(FindMovie))

	// USER

	wsUSer := &restful.WebService{}
	wsUSer.Path("/v1/users")

	wsUSer.Route(wsUSer.POST("/").
		Doc("Create new user").
		Writes("").
		Returns(201, "Created", "").
		Returns(400, "Bad request, fields not validated", model.CustomError{}.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User",
			model.ErrInternalApiUnprocessableEntity.CodeError()).
		Filter(filterUser(false)).
		To(CreateUser))

	wsUSer.Route(wsUSer.GET("/{username}").
		Doc("Get user with username").
		Param(wsUSer.PathParameter("username", "username of sought user").DataType("string")).
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(false)).
		To(GetUser))

	wsUSer.Route(wsUSer.GET("/{username}/exists").
		Doc("Know if username is already taken").
		Param(wsUSer.PathParameter("username", "username of sought user").DataType("string")).
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Returns(404, "User not found", model.ErrInternalDatabaseResourceNotFound.CodeError()).
		Filter(filterUser(false)).
		To(UsernameExists))

	return []*restful.WebService{wsRoot, wsMovie, wsUSer}
}

// Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func filterUser(checkToken bool) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		if checkToken {
			username := req.HeaderParameter(userHeaderParameter)
			logger.Sugar.Debugf("Token will be checked for this resource with username %s", username)
			if username == "" {
				res.WriteHeaderAndEntity(model.ErrInternalApiUserCredentialsNotFound.HttpCode(),
					model.ErrInternalApiUserCredentialsNotFound.CodeError())
				return
			} else if !service.UsernameAlreadyExists(username) {
				res.WriteHeaderAndEntity(model.ErrInternalApiUserBadCredentials.HttpCode(),
					model.ErrInternalApiUserBadCredentials.CodeError())
				return
			}
		}
		chain.ProcessFilter(req, res)
	}
	return filter
}