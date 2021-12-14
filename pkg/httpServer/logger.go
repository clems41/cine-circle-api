package httpServer

import (
	"github.com/emicklei/go-restful"
	"strings"
)

// Logger interface specify which functions that are needed to use a logger with this httpServer package
type Logger interface {
	Print(v ...interface{})
	Printf(template string, args ...interface{})
}

// LogRequest : Add filter for logging request (except /health/ok to avoid spamming logs)
func LogRequest(logger Logger) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		if !strings.Contains(req.Request.URL.String(), "/health/ok") {
			logger.Printf("%-10s %s", req.Request.Method, req.Request.URL.String())
		}
		chain.ProcessFilter(req, res)
	}
	return filter
}
