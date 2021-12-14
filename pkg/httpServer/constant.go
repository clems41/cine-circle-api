package httpServer

const (
	userInfoClaims  = "sub"
	tokenHeaderName = "Authorization"
	tokenKind       = "Bearer"
	tokenDelimiter  = " "
	tokenPart       = 2
)

// Environment variables and related default values used to configure HTTP server
const (
	envApplicationName              = "APPLICATION_NAME"
	envJwtRsa256PublicKey           = "JWT_RSA256_PUBLIC_KEY"
	envJwtRsa256PrivateKey          = "JWT_RSA256_PRIVATE_KEY"
	envTokenExpirationDurationHours = "TOKEN_EXPIRATION_HOURS"
	envSwaggerUrl                   = "SWAGGER_URL"
	envTracingLog                   = "HTTP_TRACING_LOG"
	envRequestLog                   = "HTTP_REQUEST_LOG"
	envBindAddress                  = "HTTP_BIND_ADDRESS"
	envReadTimeout                  = "HTTP_READ_TIMEOUT_SEC"
	envReadHeaderTimeout            = "HTTP_READ_HEADER_TIMEOUT_SEC"
	envWriteTimeout                 = "HTTP_WRITE_TIMEOUT_SEC"
	envIdleTimeout                  = "HTTP_IDLE_TIMEOUT_SEC"
	envMaxHeaderBytes               = "HTTP_MAX_HEADER_BYTES"
)

const (
	defaultSkipAuthentication           = "false"
	defaultTokenExpirationDurationHours = "24"
	defaultSwaggerUrl                   = "/swagger.json"
	defaultTracingLog                   = "false"
	defaultRequestLog                   = "true"
	defaultBindAddress                  = ":8080"
	defaultReadTimeout                  = "0"
	defaultReadHeaderTimeout            = "0"
	defaultWriteTimeout                 = "0"
	defaultIdleTimeout                  = "0"
	defaultMaxHeaderBytes               = "0"
)
