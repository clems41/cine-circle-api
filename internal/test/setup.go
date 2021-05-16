package test

import (
	"bytes"
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/postgres"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"cine-circle/internal/webService"
	"cine-circle/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var DatabaseIndexName = 0

type TestingHTTPServer struct {
	container *restful.Container
	t         *testing.T
	token     *string
}

type KeyValue struct {
	Key string
	Value string
}

func OpenDatabase(t *testing.T) (DB *gorm.DB, closeFunction func()) {
	// Open connection with real database for creating testing database
	realDB, err := postgres.OpenConnection()
	if err != nil {
		t.Fatalf(err.Error())
	}

	testingDatabaseName := fmt.Sprintf("%s%d", TestingDatabaseNamePrefix, DatabaseIndexName)
	DatabaseIndexName++

	// Deleting testing database if already exists before creating new one from template
	err = realDB.
		Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testingDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	err = realDB.
		Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", testingDatabaseName, TemplateDatabaseName)).
		Error
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Using new clean database for testing
	testingDB, err := postgres.OpenConnection(testingDatabaseName)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create closeFunction : to be used in testcases
	closeFunction = func() {
		// closing testing database before deleting it from real database
		err = postgres.CloseConnection(testingDB)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// delete testing database
		err = realDB.
			Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testingDatabaseName)).
			Error
		if err != nil {
			logger.Sugar.Fatalf(err.Error())
		}

		// closing connection from real database
		err = postgres.CloseConnection(realDB)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	webService.ActualUserHandler = webService.NewActualUserHandler(testingDB)

	return testingDB, closeFunction
}

func NewTestingHTTPServer(t *testing.T, handlers ...webService.Handler) (server *TestingHTTPServer) {
	logger.InitLogger()
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	//Create route for receiving request
	server = &TestingHTTPServer{
		container: restful.NewContainer(),
		t:         t,
	}
	for _, handler := range handlers {
		for _, ws := range handler.WebServices() {
			server.container.Add(ws)
		}
	}
	return
}

func (httpServer *TestingHTTPServer) AuthenticateUserPermanently(user *repositoryModel.User) (err error) {
	if user == nil {
		return typedErrors.NewAuthenticationErrorf("cannot authenticate user : user is nil")
	}
	if user.Username == nil {
		return typedErrors.NewAuthenticationErrorf("cannot authenticate user : username is nil for user id %d", user.GetID())
	}
	token, err := utils.GenerateTokenWithUserID(user.GetID())
	if err != nil {
		return
	}
	httpServer.token = &token

	return
}

func (httpServer *TestingHTTPServer) SendRequest(url, method string) *http.Response {
	return httpServer.sendRequest(url, method, nil, nil, nil)
}

func (httpServer *TestingHTTPServer) SendRequestWithQueryParameters(url, method string, queryParameters []KeyValue) *http.Response {
	return httpServer.sendRequest(url, method, queryParameters, nil, nil)
}

func (httpServer *TestingHTTPServer) SendRequestWithBody(url, method string, body interface{}) *http.Response {
	return httpServer.sendRequest(url, method, nil, nil, body)
}

func (httpServer *TestingHTTPServer) SendRequestWithHeaders(url, method string, headers []KeyValue) *http.Response {
	return httpServer.sendRequest(url, method, nil, headers, nil)
}

func (httpServer *TestingHTTPServer) sendRequest(url, method string, queryParameters, headerParameters []KeyValue, body interface{}) *http.Response {
	// setup request + writer
	httpWriter := httptest.NewRecorder()

	// Adding body to request if provided
	var jsonData []byte
	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			httpServer.t.Fatalf("Error occurs when marhalling data : %s", err)
		}
	}

	// Creating request with data
	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		httpServer.t.Fatalf("Error while creating request for url %s and method %s", url, method)
	}

	// Adding header for handling JSON
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	// Adding authentication header if token has been provided
	if httpServer.token != nil {
		httpRequest.Header.Set(constant.AuthenticationHeaderName, "Bearer " + *httpServer.token)
	}

	// Adding headers if provided
	for _, header := range headerParameters {
		httpRequest.Header.Set(header.Key, header.Value)
	}

	// Adding query parameters if provided
	query := httpRequest.URL.Query()
	for _, queryParameter := range queryParameters {
		query.Add(queryParameter.Key, queryParameter.Value)
	}
	httpRequest.URL.RawQuery = query.Encode()

	// Sending request
	httpServer.container.ServeHTTP(httpWriter, httpRequest)
	return httpWriter.Result()
}

func (httpServer *TestingHTTPServer) DecodeResponse(response *http.Response, out interface{}) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(httpServer.t, err)
	err = json.Unmarshal(bodyBytes, out)
	require.NoError(httpServer.t, err)
}
