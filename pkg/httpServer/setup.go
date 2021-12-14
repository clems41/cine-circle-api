package httpServer

import (
	"cine-circle-api/pkg/httpServer/swagger"
	"cine-circle-api/pkg/utils/envUtils"
	"github.com/emicklei/go-restful"
	restfulSpec "github.com/emicklei/go-restful-openapi"
	"net/http"
	"strconv"
	"time"
)

// NewRestfulContainer return pointer on RestfulContainer. This container can be used to configure and run HTTP Server.
func NewRestfulContainer() *RestfulContainer {
	container := restful.NewContainer()

	// Allow JSON in requests and responses
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	// Autorise le gzipAllow gzip encoding
	container.EnableContentEncoding(true)

	// Add filter to handle properly CORS
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      container}
	container.Filter(cors.Filter)
	container.Filter(container.OPTIONSFilter)

	return &RestfulContainer{
		container: container,
	}
}

// SetLogger specify logger to use with HTTP server logging
//  - HTTP_REQUEST_LOG : if set to true, all received requests will be logged, except /health/ok (default : true)
// 	- HTTP_TRACING_LOG : if set to true will activate tracing log (default : false)
func (rc *RestfulContainer) SetLogger(logger Logger) {
	enableTracingLog := envUtils.GetFromEnvOrDefault(envTracingLog, defaultTracingLog) == "true"
	enableRequestLog := envUtils.GetFromEnvOrDefault(envRequestLog, defaultRequestLog) == "true"

	// Can be used to add more logs
	if enableRequestLog {
		rc.container.Filter(LogRequest(logger))
	}
	if enableTracingLog {
		restful.TraceLogger(logger)
		restful.EnableTracing(true)
	}
	rc.logger = logger
}

// Container return restful container pointer
func (rc *RestfulContainer) Container() *restful.Container {
	return rc.container
}

// GenerateSwagger will register all webServices into swagger.
// You can specify swagger URL using env variable SWAGGER_URL (default: /swagger.json)
// Add endpoint at swaggerUrl for getting swagger.json file (based on documentation from each webService)
// You can visualize it here : http://swagger-ui.default.svc.kube.isi/?url=http://localhost:8080/swagger.json
// You also can specify some parameters in swagger using enrichSwaggerObject func argument.
// Example here : https://github.com/emicklei/go-restful/blob/v3/examples/openapi/restful-openapi.go#L151
func (rc *RestfulContainer) GenerateSwagger(swaggerInfo swagger.Info) {
	swaggerUrl := envUtils.GetFromEnvOrDefault(envSwaggerUrl, defaultSwaggerUrl)
	config := restfulSpec.Config{
		WebServices: rc.container.RegisteredWebServices(),
		APIPath:     swaggerUrl,
	}
	config.PostBuildSwaggerObjectHandler = swagger.EnrichSwaggerObject(swaggerInfo)
	rc.container.Add(restfulSpec.NewOpenAPIService(config))
}

// AddHandlers add webServices from handlers into restful container
func (rc *RestfulContainer) AddHandlers(handlers ...Handler) {
	for _, handler := range handlers {
		webService := handler.WebService()
		rc.container.Add(webService)
		if rc.logger != nil {
			for _, route := range webService.Routes() {
				rc.Printf("%-10s %s", route.Method, route.Path)
			}
		}
	}
	rc.Printf("\n-----------------------------------------------------------------------------------------------------------------\n\n")
}

// AddFilter add webServices from handlers into restful container
func (rc *RestfulContainer) AddFilter(filter restful.FilterFunction) {
	rc.container.Filter(filter)
}

// HttpServer will create and return http.Server from container to serve all endpoints.
// You can change settings using environment variables.
//  - SWAGGER_URL
//  - HTTP_TRACING_LOG
//  - HTTP_REQUEST_LOG
//  - HTTP_BIND_ADDRESS
//  - HTTP_READ_TIMEOUT_SEC
//  - HTTP_READ_HEADER_TIMEOUT_SEC
//  - HTTP_WRITE_TIMEOUT_SEC
//  - HTTP_IDLE_TIMEOUT_SEC
//  - HTTP_MAX_HEADER_BYTES
// You can update http.Server parameters before starting it.
func (rc *RestfulContainer) HttpServer() (server *http.Server, err error) {
	readTimeoutStr := envUtils.GetFromEnvOrDefault(envReadTimeout, defaultReadTimeout)
	readTimeout, err := strconv.Atoi(readTimeoutStr)
	if err != nil {
		return
	}
	readHeaderTimeoutStr := envUtils.GetFromEnvOrDefault(envReadHeaderTimeout, defaultReadHeaderTimeout)
	readHeaderTimeout, err := strconv.Atoi(readHeaderTimeoutStr)
	if err != nil {
		return
	}
	writeTimeoutStr := envUtils.GetFromEnvOrDefault(envWriteTimeout, defaultWriteTimeout)
	writeTimeout, err := strconv.Atoi(writeTimeoutStr)
	if err != nil {
		return
	}
	idleTimeoutStr := envUtils.GetFromEnvOrDefault(envIdleTimeout, defaultIdleTimeout)
	idleTimeout, err := strconv.Atoi(idleTimeoutStr)
	if err != nil {
		return
	}
	maxHeaderBytesStr := envUtils.GetFromEnvOrDefault(envMaxHeaderBytes, defaultMaxHeaderBytes)
	maxHeaderBytes, err := strconv.Atoi(maxHeaderBytesStr)
	if err != nil {
		return
	}
	// Create HTTP server
	server = &http.Server{
		Addr:              envUtils.GetFromEnvOrDefault(envBindAddress, defaultBindAddress),
		Handler:           rc.Container(),
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	return
}

// Printf will log using Printf method if logger has been set
func (rc *RestfulContainer) Printf(template string, args ...interface{}) {
	if rc.logger != nil {
		rc.logger.Printf(template, args...)
	}
}
