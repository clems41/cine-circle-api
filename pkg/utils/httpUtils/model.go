package httpUtils

import (
	"fmt"
	"github.com/emicklei/go-restful"
)

type Request struct {
	Url               string
	Method            string
	ContentType       string
	QueryParameters   map[string]string
	HeadersParameters map[string]string
	Body              interface{}
}

type Parameter struct {
	Name         string
	Description  string
	DefaultValue interface{}
	DataType     string
	Required     bool
}

func (path Parameter) PathParameter() *restful.Parameter {
	return restful.
		PathParameter(path.Name, path.Description).
		DataType(path.DataType).
		Required(path.Required).
		DefaultValue(fmt.Sprintf("%v", path.DefaultValue))
}

func (path Parameter) QueryParameter() *restful.Parameter {
	return restful.
		QueryParameter(path.Name, path.Description).
		DataType(path.DataType).
		Required(path.Required).
		DefaultValue(fmt.Sprintf("%v", path.DefaultValue))
}

func (path Parameter) Joker() string {
	return fmt.Sprintf("{%s}", path.Name)
}
