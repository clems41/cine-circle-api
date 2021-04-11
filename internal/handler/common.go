package handler

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
)

type handler interface {
	WebService() *restful.WebService
}

func AddWebService(container *restful.Container, handler handler) {
	container.Add(handler.WebService())
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