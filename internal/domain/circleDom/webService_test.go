package circleDom

import (
	"cine-circle-api/internal/repository/instance/circleRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/internal/test/setupTestCase"
	"cine-circle-api/internal/test/testSampler"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"strings"
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
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with wrong form (missing name) --> NOK 400
	wrongForm := correctForm
	wrongForm.Name = ""
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct form --> OK 201
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check view fields
	var view CreateView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":          testRuler.NotEmptyField{},
			"Name":        correctForm.Name,
			"Description": correctForm.Description,
			"Users":       testRuler.EmptyField{},
		},
	})

	// Check that circle has been created into database
	var circle model.Circle
	err := db.
		Take(&circle, view.Id).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, circle.ID)
}

/* Update */

func TestHandler_Update(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	existingCircle := sampler.GetCircle()
	correctForm := UpdateForm{
		CommonForm: CommonForm{
			Name:        fake.Title(),
			Description: fake.Sentences(),
		},
	}
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingCircle.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(correctTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithBody(wrongTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrong form (missing name) --> NOK 400
	wrongForm := correctForm
	wrongForm.Name = ""
	resp = httpMock.SendRequestWithBody(correctTestPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct form --> OK 200
	resp = httpMock.SendRequestWithBody(correctTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view UpdateView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":          existingCircle.ID,
			"Name":        correctForm.Name,
			"Description": correctForm.Description,
			"Users":       testRuler.EmptyField{},
		},
	})

	// Check that circle has been updated from database
	err := db.
		Take(&existingCircle, existingCircle.ID).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, existingCircle.ID)
}

/* Delete */

func TestHandler_Delete(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	existingCircle := sampler.GetCircle()
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingCircle.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodDelete)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correct circleId --> NOK 200
	resp = httpMock.SendRequest(correctTestPath, http.MethodDelete)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that circle has been deleted into database
	err := db.
		Take(&existingCircle, existingCircle.ID).
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
	existingCircle := sampler.GetCircle()
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingCircle.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correct circleId --> NOK 200
	resp = httpMock.SendRequest(correctTestPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view GetView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":          existingCircle.ID,
			"Name":        existingCircle.Name,
			"Description": existingCircle.Description,
			"Users":       testRuler.EmptyField{},
		},
	})
}

/* Search */

func TestHandler_Search(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, _, tearDown := setupTestcase(t, false)
	defer tearDown()

	// Create testing data
	keyword := fake.Words() + fake.Words()
	var circles []*model.Circle
	nbCirclesNameMatching := 6
	for idx := range fakeData.FakeRange(12, 18) { // create between 12 and 18 circles (random number) but only 6 with matching name
		if idx < nbCirclesNameMatching {
			circles = append(circles, sampler.GetCircleWithName(fake.Words()+keyword+fake.Word()))
		} else {
			circles = append(circles, sampler.GetCircle())
		}
	}
	user := sampler.GetUser()
	page := 1
	pageSize := 10 // to be sure that all circles will fit on 2 pages (10 < 12 < 18 < 2*10)
	keywords := map[string]string{
		circleNameQueryParameter.Name: keyword,
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, nil, keywords))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with name filter using keyword
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, nil, keywords))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if view is correctly sorted
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbCirclesNameMatching, view.NumberOfItems)
	require.Equal(t, nbCirclesNameMatching, len(view.Circles))
	for _, viewCircle := range view.Circles {
		require.True(t, strings.Contains(strings.ToLower(viewCircle.Name), strings.ToLower(keyword)),
			"Circle name %s should contains keyword %s", strings.ToLower(viewCircle.Name), strings.ToLower(keyword))
	}
}

/* Add user */

func TestHandler_AddUser(t *testing.T) {

}

/* Delete user */

func TestHandler_DeleteUser(t *testing.T) {

}

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = setupTestCase.OpenCleanDatabaseFromTemplate(t)
	repo := circleRepository.New(db)
	svc := NewService(repo)
	ws := NewHandler(svc)
	httpMock = httpServerMock.New(t, logger.Logger(), ws)
	sampler = testSampler.New(t, db, populateDatabase)
	ruler = testRuler.New(t)
	return
}

// searchQueryParameters will return queryParameters based on search parameters
func searchQueryParameters(page, pageSize int, sort []string, keywords map[string]string) (queryParameters map[string][]string) {
	queryParameters = map[string][]string{
		pageQueryParameter.Name:     {fmt.Sprintf("%d", page)},
		pageSizeQueryParameter.Name: {fmt.Sprintf("%d", pageSize)},
	}
	for _, sortValue := range sort {
		queryParameters["sort"] = append(queryParameters["sort"], sortValue)
	}
	for keywordField, keywordValue := range keywords {
		queryParameters[keywordField] = []string{keywordValue}
	}
	return
}
