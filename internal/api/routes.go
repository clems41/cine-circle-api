package api

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
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
	wsMovie.Path("/v1/movie")

	wsMovie.Route(wsMovie.GET("/").
		Doc("Get movie by ID or by title").
		Param(wsMovie.QueryParameter("id", "Get movie by ID (based on IMDb ids)").DataType("string")).
		Param(wsMovie.QueryParameter("title", "Get movie by title").DataType("string")).
		Writes(model.Movie{}).
		Returns(200, "OK", model.Movie{}).
		Filter(filterUser()).
		To(FindMovie))

	// USER

	wsUSer := &restful.WebService{}
	wsUSer.Path("/v1/user")

	wsUSer.Route(wsUSer.POST("/").
		Doc("Create new user").
		Writes("").
		Returns(201, "Created", "").
		Returns(400, "Bad request, fields not validated", model.CustomError{}.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to User", model.CustomError{}.CodeError()).
		Filter(filterUser()).
		To(CreateUser))

	wsUSer.Route(wsUSer.GET("/").
		Doc("Get user by ID").
		Param(wsUSer.QueryParameter("id", "Get user by ID").DataType("string")).
		Writes(model.User{}).
		Returns(200, "OK", model.User{}).
		Filter(filterUser()).
		To(GetUser))

	return []*restful.WebService{wsRoot, wsMovie, wsUSer}
}

// Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func filterUser() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}