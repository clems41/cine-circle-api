package httpServerMock

import (
	"bytes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// SendRequest will send request without any header parameters, query parameters or body then return related HTTP response
func (server *Server) SendRequest(url, method string) *http.Response {
	return server.sendRequest(url, method, nil, nil, nil)
}

// SendRequestWithQueryParameters will send request with specific query parameters using map[string]string (key & value) then return related HTTP response
func (server *Server) SendRequestWithQueryParameters(url, method string, queryParameters map[string][]string) *http.Response {
	return server.sendRequest(url, method, queryParameters, nil, nil)
}

// SendRequestWithBody will send request with JSON body based on interface argument then return related HTTP response
func (server *Server) SendRequestWithBody(url, method string, body interface{}) *http.Response {
	return server.sendRequest(url, method, nil, nil, body)
}

// SendRequestWithHeaderParameters will send request with headers parameters using map[string]string (key & value) then return related HTTP response
func (server *Server) SendRequestWithHeaderParameters(url, method string, headers map[string]string) *http.Response {
	return server.sendRequest(url, method, nil, headers, nil)
}

// sendRequest is generic method for sending request with bogy, query parameters and/or header parameters then return related HTTP response
func (server *Server) sendRequest(url, method string, queryParameters map[string][]string, headerParameters map[string]string, body interface{}) *http.Response {
	// setup request + writer
	httpWriter := httptest.NewRecorder()

	// Adding body to request if provided
	var jsonData []byte
	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			server.t.Fatalf("Error occurs when marhalling data : %s", err)
		}
	}

	// Creating request with data
	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		server.t.Fatalf("Error while creating request for url %s and method %s", url, method)
	}

	// Adding header for handling JSON
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)

	// Adding authentication header if token has been provided
	if server.token != nil {
		httpRequest.Header.Set(tokenHeaderName,  tokenKind+tokenDelimiter+ *server.token)
	}

	// Adding headers if provided
	for headerKey, headerValue := range headerParameters {
		httpRequest.Header.Set(headerKey, headerValue)
	}

	// Adding query parameters if provided
	query := httpRequest.URL.Query()
	for queryKey, queryValues := range queryParameters {
		for _, queryValue := range queryValues {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	// Sending request
	server.container.Container().ServeHTTP(httpWriter, httpRequest)
	return httpWriter.Result()
}

// DecodeResponse will unmarshal response to out interface (need pointer)
func (server *Server) DecodeResponse(response *http.Response, out interface{}) {
	defer func() {
		err := response.Body.Close()
		require.NoError(server.t, err)
	}()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(server.t, err)
	err = json.Unmarshal(bodyBytes, out)
	require.NoError(server.t, err)
}
