package libraryDom

import "github.com/emicklei/go-restful"

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (api handler) WebServices() (webServices []*restful.WebService) {
	wsLibrary := &restful.WebService{}
	webServices = append(webServices, wsLibrary)

	wsLibrary.Path("/v1/library")

	return
}
