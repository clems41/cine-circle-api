// Package httpServer can be used to create HTTP Server.
//  This package also define method to generate token with user information as claims and also a method to get user information from token string.
//  NewRestfulContainer create container that will be used to define all your routes, which logger to use, swagger URL, etc...
//  From this container, you can obviously start a new HTTP server using some parameters through environment variables :
//  - SWAGGER_URL : url to get swagger.json file (useful to generate swagger from SwaggerUI (default : /swagger.json)
//  - HTTP_TRACING_LOG : get more info from request and trace all requests (default : false)
//  - HTTP_REQUEST_LOG : log all received request except liveness probe through /health/ok endpoint (default : true)
//  - HTTP_BIND_ADDRESS : define which port to use to serve HTTP server (default : :8080)
//  - HTTP_READ_TIMEOUT_SEC : maximum duration for reading the entire request, including the body (default : 0)
//  - HTTP_READ_HEADER_TIMEOUT_SEC : amount of time allowed to read request headers (default : 0)
//  - HTTP_WRITE_TIMEOUT_SEC : maximum duration before timing out writes of the response (default : 0)
//  - HTTP_IDLE_TIMEOUT_SEC : maximum amount of time to wait for the next request when keep-alives are enabled (default : 0)
//  - HTTP_MAX_HEADER_BYTES : maximum number of bytes the server will read parsing the request header's keys and values, including the request line (default : 0 -> 1MB)
//  - TOKEN_EXPIRATION_HOURS : duration before token will expire (default : 24 -> 1 day)
package httpServer
