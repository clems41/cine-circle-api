package httpServerMock

import (
	"cine-circle-api/pkg/httpServer"
	"os"
	"testing"
)

type Server struct {
	container *httpServer.RestfulContainer
	t         *testing.T
	token     *string
}

// New return new Server that can be used as mock in tests
// To work, you need to provide which handlers you want to test
// You can also define which logger to use. If nil, no logs will be printed from mock.
func New(t *testing.T, logger httpServer.Logger, handlers ...httpServer.Handler) (mock *Server) {
	//Create route for receiving request
	mock = &Server{
		container: httpServer.NewRestfulContainer(),
		t:         t,
	}
	if logger != nil {
		mock.container.SetLogger(logger)
	}
	mock.container.AddHandlers(handlers...)
	mock.setJwtRsaKeys()
	return
}

// BasePath return path of specified webService (index)
func (server *Server) BasePath(index int) string {
	webServices := server.container.Container().RegisteredWebServices()
	if index >= len(webServices) {
		server.t.Fatalf("Cannot access webService index %d, httpServer got only %d webServices", index, len(webServices))
	}
	return webServices[index].RootPath()
}

// AuthenticateUserPermanently authenticate user specified in parameter.
// Token will be generated and added to all next requests during test execution.
func (server *Server) AuthenticateUserPermanently(user interface{}) {
	if user == nil {
		server.t.Fatalf("cannot authenticate user : user is nil")
	}
	token, err := httpServer.GenerateTokenWithUserInfo(user)
	if err != nil {
		server.t.Fatalf(err.Error())
	}
	server.token = &token.TokenString

	return
}

func (server *Server) setJwtRsaKeys() {
	err := os.Setenv(envJwtRsa256PrivateKey, jwtRsa256PrivateKey)
	if err != nil {
		server.t.Fatalf(err.Error())
	}
	err = os.Setenv(envJwtRsa256PublicKey, jwtRsa256PublicKey)
	if err != nil {
		server.t.Fatalf(err.Error())
	}
}
