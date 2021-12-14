package httpUtils

import "fmt"

type Request struct {
	Url               string
	Method            string
	ContentType       string
	QueryParameters   map[string]string
	HeadersParameters map[string]string
	Body              interface{}
}

type PathParameter string

func (path PathParameter) Joker() string {
	return fmt.Sprintf("{%s}", path)
}

func (path PathParameter) String() string {
	return string(path)
}
