package healthCheckDom

import (
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	handlerToTest := NewHandler()
	testingHTTPServer := httpServerMock.New(t, logger.Logger(), handlerToTest)

	// Routes use for this test
	basePath := testingHTTPServer.BasePath(0)

	// Send request and check response code, should return 200
	resp := testingHTTPServer.SendRequest(basePath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
