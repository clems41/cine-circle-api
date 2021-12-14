package healthCheckDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"net/http"
)

var _ Handler = (*handler)(nil)

type Handler interface {
}

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

// WebService is only useful to Kubernetes. When application will be running in Pod, regularly http request will be sent to /health/ok in order to check server status.
// If server doesn't respond with 200 OK, Kubenertes will update Pod status as not running.
func (hd *handler) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)

	ws.Path("/health/ok")
	tags := []string{swaggerConst.OtherTag}

	ws.Route(ws.GET("/").
		Produces(restful.MIME_JSON).
		Doc("Simple API health check").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(func(req *restful.Request, res *restful.Response) {
			res.WriteHeader(http.StatusOK)
		}))

	return
}
