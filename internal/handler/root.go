package handler

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func NewRootHandler() *restful.WebService {
	wsRoot := &restful.WebService{}

	// HEALTH
	wsRoot.Route(wsRoot.GET("/health/ok").
		Doc("Simple API health check").
		To(func(req *restful.Request, res *restful.Response) {
			res.WriteHeader(http.StatusOK)
		}))
	return wsRoot
}
