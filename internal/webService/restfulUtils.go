package webService

import (
	"cine-circle/pkg/logger"
	"github.com/emicklei/go-restful"
)

type Handler interface {
	WebServices() []*restful.WebService
}

// AddHandlersToRestfulContainer : add webServices into restful container from handlers
func AddHandlersToRestfulContainer(container *restful.Container, handlers ...Handler) {
	for _, handler := range handlers {
		webServices := handler.WebServices()
		for _, webService := range webServices {
			container.Add(webService)
			for _, route := range webService.Routes() {
				logger.Sugar.Infof("%s \t %s", route.Method, route.Path)
			}
		}
	}
}
