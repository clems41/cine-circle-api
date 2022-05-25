package circleDom

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/model/testSampler"
	"cine-circle-api/internal/repository"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/sql/sqlTest"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"fmt"
	"github.com/icrowley/fake"
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
			"Users":       testRuler.NotEmptyField{},
		},
	})
	require.Len(t, view.Users, 1) // user that created circle should be automatically added into it

	// Check that circle has been created into database
	var circle model.Circle
	err := db.
		Preload("Users").
		Take(&circle, view.Id).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, circle.ID)
	require.Len(t, circle.Users, 1) // user that created circle should be automatically added into it
}

/* Update */

func TestHandler_Update(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	userNotFromCircle := sampler.GetUser()
	existingCircle := sampler.GetCircleWithUsers()
	userFromCircle := existingCircle.Users[0]
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

	// Try with user not from circle --> NOK 404
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequestWithBody(wrongTestPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(userFromCircle)
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
			"Users":       testRuler.NotEmptyField{},
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
	userNotFromCircle := sampler.GetUser()
	existingCircle := sampler.GetCircleWithUsers()
	userFromCircle := existingCircle.Users[0]
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingCircle.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodDelete)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user not from circle --> NOK 404
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(userFromCircle)
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
	userNotFromCircle := sampler.GetUser()
	existingCircle := sampler.GetCircleWithUsers()
	userFromCircle := existingCircle.Users[0]
	wrongTestPath := fmt.Sprintf("%s/%d", testPath, 9865)
	correctTestPath := fmt.Sprintf("%s/%d", testPath, existingCircle.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctTestPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user not from circle--> NOK 404
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequest(wrongTestPath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with non-existing circleId --> NOK 404
	httpMock.AuthenticateUserPermanently(userFromCircle)
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
			"Users":       testRuler.NotEmptyField{},
		},
	})
}

/* Search */

// TestHandler_Search check that circles returned are only one with authenticated user in it
func TestHandler_Search(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	userNotFromCircle := sampler.GetUser()
	userFromCircle := sampler.GetUser()
	var circles []*model.Circle
	nbCirclesNameMatching := 6
	for idx := range fakeData.FakeRange(12, 18) { // create between 12 and 18 circles (random number) but only 6 with matching name
		if idx < nbCirclesNameMatching {
			circles = append(circles, sampler.GetCircleWithSpecificUser(userFromCircle))
		} else {
			circles = append(circles, sampler.GetCircleWithUsers())
		}
	}
	page := 1
	pageSize := 10

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, nil, nil))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user with no circles
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, nil, nil))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view that contains no circles
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 0, view.NumberOfPages)
	require.Equal(t, 0, view.NumberOfItems)
	require.Equal(t, 0, len(view.Circles))

	// Try with user with no circles
	httpMock.AuthenticateUserPermanently(userFromCircle)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, nil, nil))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view that contains 6 circles
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbCirclesNameMatching, view.NumberOfItems)
	require.Equal(t, nbCirclesNameMatching, len(view.Circles))
}

/* Add user */

func TestHandler_AddUser(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	userNotFromCircle := sampler.GetUser()
	existingCircle := sampler.GetCircleWithUsers()
	userFromCircle := existingCircle.Users[0]
	userToAdd := sampler.GetUser()
	wrongCircleId := 999
	wrongUserId := 985
	wrongCircleWrongUserPath := fmt.Sprintf("%s/%d/%d", testPath, wrongCircleId, wrongUserId)
	wrongCircleCorrectUserPath := fmt.Sprintf("%s/%d/%d", testPath, wrongCircleId, userToAdd.ID)
	correctCircleWrongUserPath := fmt.Sprintf("%s/%d/%d", testPath, existingCircle.ID, wrongUserId)
	correctCircleCorrectUserPath := fmt.Sprintf("%s/%d/%d", testPath, existingCircle.ID, userToAdd.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctCircleCorrectUserPath, http.MethodPut)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user not from circle --> NOK 404
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequest(wrongCircleWrongUserPath, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrongCircleWrongUserPath --> NOK 404
	httpMock.AuthenticateUserPermanently(userFromCircle)
	resp = httpMock.SendRequest(wrongCircleWrongUserPath, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrongCircleCorrectUserPath --> NOK 404
	resp = httpMock.SendRequest(wrongCircleCorrectUserPath, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correctCircleWrongUserPath --> NOK 404
	resp = httpMock.SendRequest(correctCircleWrongUserPath, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correctCircleCorrectUserPath --> OK 200
	resp = httpMock.SendRequest(correctCircleCorrectUserPath, http.MethodPut)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view AddUserView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":          existingCircle.ID,
			"Name":        existingCircle.Name,
			"Description": existingCircle.Description,
			"Users":       testRuler.NotEmptyField{},
		},
	})

	// Check that user has been added into circle in database
	err := db.
		Preload("Users").
		Take(existingCircle).
		Error
	require.NoError(t, err)
	var userFound bool
	for _, circleUser := range existingCircle.Users {
		if circleUser.ID == userToAdd.ID {
			userFound = true
		}
	}
	require.True(t, userFound, "User should be found in circle")
}

/* Delete user */

func TestHandler_DeleteUser(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	userNotFromCircle := sampler.GetUser()
	existingCircle := sampler.GetCircleWithUsers()
	userFromCircle := existingCircle.Users[0]
	userToDelete := sampler.GetUser()
	wrongCircleId := 999
	wrongUserId := 985
	wrongCircleWrongUserPath := fmt.Sprintf("%s/%d/%d", testPath, wrongCircleId, wrongUserId)
	wrongCircleCorrectUserPath := fmt.Sprintf("%s/%d/%d", testPath, wrongCircleId, userToDelete.ID)
	correctCircleWrongUserPath := fmt.Sprintf("%s/%d/%d", testPath, existingCircle.ID, wrongUserId)
	correctCircleCorrectUserPath := fmt.Sprintf("%s/%d/%d", testPath, existingCircle.ID, userToDelete.ID)

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctCircleCorrectUserPath, http.MethodDelete)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user not from circle --> NOK 404
	httpMock.AuthenticateUserPermanently(userNotFromCircle)
	resp = httpMock.SendRequest(wrongCircleWrongUserPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrongCircleWrongUserPath --> NOK 404
	httpMock.AuthenticateUserPermanently(userFromCircle)
	resp = httpMock.SendRequest(wrongCircleWrongUserPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrongCircleCorrectUserPath --> NOK 404
	resp = httpMock.SendRequest(wrongCircleCorrectUserPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correctCircleWrongUserPath --> NOK 404
	resp = httpMock.SendRequest(correctCircleWrongUserPath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with correctCircleCorrectUserPath --> OK 200
	resp = httpMock.SendRequest(correctCircleCorrectUserPath, http.MethodDelete)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view DeleteUserView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":          existingCircle.ID,
			"Name":        existingCircle.Name,
			"Description": existingCircle.Description,
			"Users":       testRuler.NotEmptyField{},
		},
	})

	// Check that user has been added into circle in database
	err := db.
		Preload("Users").
		Take(existingCircle).
		Error
	require.NoError(t, err)
	var userFound bool
	for _, circleUser := range existingCircle.Users {
		if circleUser.ID == userToDelete.ID {
			userFound = true
		}
	}
	require.Falsef(t, userFound, "User should not be found in circle")
}

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = sqlTest.OpenCleanDatabaseFromTemplate(t)
	repo := repository.New(db)
	userRepo := repository.New(db)
	svc := NewService(repo, userRepo)
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
