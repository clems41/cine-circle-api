package exempleDom

import (
	"cine-circle-api/internal/repository/instance/exempleRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/internal/test/setupTestCase"
	"cine-circle-api/internal/test/testSampler"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

/* Create */

func TestHandler_Create(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	correctForm := CreateForm{CommonForm{
		// TODO add your custom fields here
	}}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// TODO add your custom testcase scenarii depending on your custom fields rules

	// Try with correct form --> OK 201
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check view fields
	var view CreateView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id": testRuler.NotEmptyField{},
			// TODO add your custom fields check here
		},
	})

	// Check that exemple has been created into database
	var exemple model.Exemple
	err := db.
		Take(&exemple, view.Id).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, exemple.ID)
	// TODO add your custom fields check here
}

/* Update */

func TestHandler_Update(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	existingExemple := sampler.GetExemple()
	correctForm := UpdateForm{
		CommonForm: CommonForm{
			// TODO add your custom fields here
		},
	}
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingExemple.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(correctTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing exempleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithBody(wrongTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// TODO add your custom testcase scenarii depending on your custom fields rules

	// Try with correct form --> OK 200
	resp = httpMock.SendRequestWithBody(correctTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view UpdateView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id": existingExemple.ID,
			// TODO add your custom fields check here
		},
	})

	// Check that exemple has been updated from database
	err := db.
		Take(&existingExemple, existingExemple.ID).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, existingExemple.ID)
	// TODO add your custom fields check here
}

/* Delete */

func TestHandler_Delete(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	existingExemple := sampler.GetExemple()
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingExemple.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodDelete)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing exempleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correct exempleId --> NOK 200
	resp = httpMock.SendRequest(correctTestPath, http.MethodDelete)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that exemple has been deleted into database
	err := db.
		Take(&existingExemple, existingExemple.ID).
		Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

/* Get */

func TestHandler_Get(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	existingExemple := sampler.GetExemple()
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingExemple.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing exempleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correct exempleId --> NOK 200
	resp = httpMock.SendRequest(correctTestPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view GetView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id": existingExemple.ID,
			// TODO add your custom fields check here
		},
	})
}

/* Search */

func TestHandler_Search_Sort(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, ruler, tearDown := setupTestcase(t, false)
	defer tearDown()

	// Create testing data
	var exemples []*model.Exemple
	for range fakeData.FakeRange(6, 12) { // create between 6 and 12 exemples (random number)
		exemples = append(exemples, sampler.GetExemple())
	}
	user := sampler.GetUser()
	page := 1
	pageSize := 15 // to be sure that all exemples can fit in one page (15 > 12)
	sort := []string{}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, sort, nil))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with default sorting (id:asc) adn default pagination (page = 1 and pageSize = 20) --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if view is correctly sorted
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, defaultPage, view.CurrentPage)
	require.Equal(t, defaultPageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, len(exemples), view.NumberOfItems)
	require.Equal(t, len(exemples), len(view.Exemples))
	ruler.SliceIsSorted(view.Exemples, func(i, j int) bool {
		return view.Exemples[i].Id < view.Exemples[j].Id
	})

	// TODO add your custom testcase scenarii depending on which sort value you want to test.
	// You can use previous example
}

func TestHandler_Search_Page2(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, _, tearDown := setupTestcase(t, false)
	defer tearDown()

	// Create testing data
	var exemples []*model.Exemple
	for range fakeData.FakeRange(12, 18) { // create between 12 and 18 exemples (random number)
		exemples = append(exemples, sampler.GetExemple())
	}
	user := sampler.GetUser()
	page := 2
	pageSize := 10 // to be sure that all exemples will fit on 2 pages (10 < 12 < 18 < 2*10)
	sort := []string{}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, sort, nil))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with sorting firstName:asc --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, sort, nil))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if view is correctly sorted
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 2, view.NumberOfPages)
	require.Equal(t, len(exemples), view.NumberOfItems)
	require.Equal(t, len(exemples)-pageSize, len(view.Exemples))
}

// TODO add your custom testcase scenarii for filtering based on your keyword fields (cf. userDom example)

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = setupTestCase.OpenCleanDatabaseFromTemplate(t)
	repo := exempleRepository.New(db)
	svc := NewService(repo)
	ws := NewHandler(svc)
	httpMock = httpServerMock.New(t, logger.Logger(), ws)
	sampler = testSampler.New(t, db, populateDatabase)
	ruler = testRuler.New(t)
	return
}

// basicAuth return Authorization header with Basic Authentication header (ex: Basic bG9naW46cGFzc3dvcmQK)
func basicAuthHeader(exemplename, password string) map[string]string {
	auth := exemplename + ":" + password
	return map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))}
}

// searchQueryParameters will return queryParameters based on search parameters
func searchQueryParameters(page, pageSize int, sort []string, keywords map[string]string) (queryParameters map[string][]string) {
	queryParameters = map[string][]string{
		"page":     {fmt.Sprintf("%d", page)},
		"pageSize": {fmt.Sprintf("%d", pageSize)},
	}
	for _, sortValue := range sort {
		queryParameters["sort"] = append(queryParameters["sort"], sortValue)
	}
	for keywordField, keywordValue := range keywords {
		queryParameters[keywordField] = []string{keywordValue}
	}
	return
}
