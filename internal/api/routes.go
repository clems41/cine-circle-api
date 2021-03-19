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

	return []*restful.WebService{wsRoot, wsMovie}
}

// Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func filterUser() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}