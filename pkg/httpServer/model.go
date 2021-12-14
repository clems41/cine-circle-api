package httpServer

import (
	"github.com/emicklei/go-restful"
)

type Handler interface {
	WebService() *restful.WebService
}

type RestfulContainer struct {
	container *restful.Container
	logger    Logger
}