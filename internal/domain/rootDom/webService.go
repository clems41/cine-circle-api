package rootDom

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (api handler) WebServices() (handlers []*restful.WebService) {
	wsRoot := &restful.WebService{}
	handlers = append(handlers, wsRoot)

	// HEALTH
	wsRoot.Route(wsRoot.GET("/health/ok").
		Doc("Simple API health check").
		To(func(req *restful.Request, res *restful.Response) {
			res.WriteHeader(http.StatusOK)
		}))
	return
}
