package circleDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/test"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing user to database
	userSample := sampler.GetUser()

	// fields for circle
	name := fake.Title()
	description := fake.Sentences()

	var creation Creation

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with all missing fields, should fail and return 400
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with missing field : name, should fail and return 400
	creation.Description = description
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with missing field : description, should fail and return 400
	creation.Description = ""
	creation.Name = name
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with correct fields, should return 201
	creation.Description = description
	creation.Name = name
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"CircleID":    test.NotEmptyField{},
		"Name":        creation.Name,
		"Description": creation.Description,
		"Users": []UserView{
			{
				UserID:      userSample.GetID(),
				Username:    *userSample.Username,
				DisplayName: userSample.DisplayName,
			},
		},
	})

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", view.CircleID).
		Error
	require.NoError(t, err)
	require.Len(t, circle.Users, 1)
}

func TestHandler_Update(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing circle to database
	circleSample := sampler.GetCircle()

	wrongBasePath := webServicePath + "/9999"
	correctBasePath := webServicePath + "/" + circleSample.GetIDAsString()

	// Add existing user to database
	userSample := sampler.GetUser()

	// fields for circle
	name := fake.Title()
	description := fake.Sentences()

	update := Update{
		Name:        name,
		Description: description,
	}

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated not in circle, should return 401
	resp = testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user in circle)
	testingHTTPServer.AuthenticateUserPermanently(&circleSample.Users[1])

	// Send request with wrong path id, should fail and return 404
	resp = testingHTTPServer.SendRequestWithBody(wrongBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with all missing field, should fail and return 400
	update = Update{
		Name:        "",
		Description: "",
	}
	resp = testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with missing field : name, should fail and return 400
	update = Update{
		Name:        "",
		Description: description,
	}
	resp = testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with missing field : description, should fail and return 400
	update = Update{
		Name:        name,
		Description: "",
	}
	resp = testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request with correct fields, should return 200
	update = Update{
		Name:        name,
		Description: description,
	}
	resp = testingHTTPServer.SendRequestWithBody(correctBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"CircleID":    circleSample.ID,
		"Name":        update.Name,
		"Description": update.Description,
		"Users":       test.NotEmptyField{},
	})
	require.Len(t, view.Users, len(circleSample.Users))
	for idx, userView := range view.Users {
		ruler.CheckStruct(userView, map[string]interface{}{
			"UserID":      circleSample.Users[idx].ID,
			"Username":    *circleSample.Users[idx].Username,
			"DisplayName": circleSample.Users[idx].DisplayName,
		})
	}

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", view.CircleID).
		Error
	require.NoError(t, err)
	require.Len(t, circle.Users, len(circleSample.Users))
}

func TestHandler_Delete(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing circle to database
	circleSample := sampler.GetCircle()

	wrongBasePath := webServicePath + "/9999"
	correctBasePath := webServicePath + "/" + circleSample.GetIDAsString()

	// Add existing user to database
	userSample := sampler.GetUser()

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated not in circle, should return 401
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user in circle)
	testingHTTPServer.AuthenticateUserPermanently(&circleSample.Users[1])

	// Send request with wrong path id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with correct path, should return 204
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", circleSample.GetID()).
		Error
	require.Error(t, err)
}

func TestHandler_Get(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing circle to database
	circleSample := sampler.GetCircle()

	wrongBasePath := webServicePath + "/9999"
	correctBasePath := webServicePath + "/" + circleSample.GetIDAsString()

	// Add existing user to database
	userSample := sampler.GetUser()

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated not in circle, should return 401
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user in circle)
	testingHTTPServer.AuthenticateUserPermanently(&circleSample.Users[1])

	// Send request with wrong path id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with all correct fields, should return 302
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"CircleID":    circleSample.ID,
		"Name":        circleSample.Name,
		"Description": circleSample.Description,
		"Users":       test.NotEmptyField{},
	})
	require.Len(t, view.Users, len(circleSample.Users))
	for idx, userView := range view.Users {
		ruler.CheckStruct(userView, map[string]interface{}{
			"UserID":      circleSample.Users[idx].ID,
			"Username":    *circleSample.Users[idx].Username,
			"DisplayName": circleSample.Users[idx].DisplayName,
		})
	}

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", view.CircleID).
		Error
	require.NoError(t, err)
	require.Len(t, circle.Users, len(circleSample.Users))
}

func TestHandler_AddUser(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing circle to database
	circleSample := sampler.GetCircle()
	// Add existing user to database
	userSample := sampler.GetUser()

	wrongBasePathWithWrongUserID := webServicePath + "/9999/" + userSample.GetIDAsString()
	wrongBasePathWithWrongCircleID := webServicePath + "/" + circleSample.GetIDAsString() + "/9999"
	wrongBasePathWithBothWrong := webServicePath + "/9999/9999"
	correctBasePath := webServicePath + "/" + circleSample.GetIDAsString() + "/" + userSample.GetIDAsString()

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodPut)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated not in circle, should return 401
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodPut)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user in circle)
	testingHTTPServer.AuthenticateUserPermanently(&circleSample.Users[1])

	// Send request with wrong path id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePathWithWrongUserID, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp = testingHTTPServer.SendRequest(wrongBasePathWithWrongCircleID, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp = testingHTTPServer.SendRequest(wrongBasePathWithBothWrong, http.MethodPut)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with all correct fields, should return 200
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodPut)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"CircleID":    circleSample.ID,
		"Name":        circleSample.Name,
		"Description": circleSample.Description,
		"Users":       test.NotEmptyField{},
	})
	// Check if user has been correctly added
	require.Len(t, view.Users, len(circleSample.Users)+1)
	for idx, userView := range view.Users {
		if idx < len(circleSample.Users) {
			ruler.CheckStruct(userView, map[string]interface{}{
				"UserID":      circleSample.Users[idx].ID,
				"Username":    *circleSample.Users[idx].Username,
				"DisplayName": circleSample.Users[idx].DisplayName,
			})
		} else {
			ruler.CheckStruct(userView, map[string]interface{}{
				"UserID":      userSample.ID,
				"Username":    *userSample.Username,
				"DisplayName": userSample.DisplayName,
			})
		}
	}

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", view.CircleID).
		Error
	require.NoError(t, err)
	require.Len(t, circle.Users, len(circleSample.Users)+1)
}

func TestHandler_DeleteUser(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Add existing circle to database
	circleSample := sampler.GetCircle()
	userFromCircleSampleToDelete := circleSample.Users[1]
	userFromCircleSample := circleSample.Users[2]
	userSample := sampler.GetUser()

	wrongBasePathWithWrongUserID := webServicePath + "/9999/" + userFromCircleSampleToDelete.GetIDAsString()
	wrongBasePathWithWrongCircleID := webServicePath + "/" + circleSample.GetIDAsString() + "/9999"
	wrongBasePathWithBothWrong := webServicePath + "/9999/9999"
	correctBasePath := webServicePath + "/" + circleSample.GetIDAsString() + "/" + userFromCircleSampleToDelete.GetIDAsString()

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated not in circle, should return 401
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user in circle)
	testingHTTPServer.AuthenticateUserPermanently(&userFromCircleSample)

	// Send request with wrong path id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePathWithWrongUserID, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp = testingHTTPServer.SendRequest(wrongBasePathWithWrongCircleID, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp = testingHTTPServer.SendRequest(wrongBasePathWithBothWrong, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with all correct fields, should return 200
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"CircleID":    circleSample.ID,
		"Name":        circleSample.Name,
		"Description": circleSample.Description,
		"Users":       test.NotEmptyField{},
	})
	// Check if user has been correctly deleted
	require.Len(t, view.Users, len(circleSample.Users)-1)

	// Check if all users has been correctly saved
	var circle repositoryModel.Circle
	err := DB.
		Preload("Users").
		Take(&circle, "id = ?", view.CircleID).
		Error
	require.NoError(t, err)
	require.Len(t, circle.Users, len(circleSample.Users)-1)
}

func TestHandler_List(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	circleWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := circleWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, circleWebService)

	// Create circles with user in it
	userSample := sampler.GetUser()
	var circles []*repositoryModel.Circle
	for range test.FakeRange(5, 15) {
		circles = append(circles, sampler.GetCircle(*userSample))
	}


	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(webServicePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request with user authenticated, should work
	queryParams := []test.KeyValue{
		{
			Key:   "page",
			Value: "1",
		},
		{
			Key:   "pageSize",
			Value: "10",
		},
	}
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParams)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var list ListView
	testingHTTPServer.DecodeResponse(resp, &list)

	require.Equal(t, len(circles), list.NumberOfItems)
	require.NotEqual(t, 0, len(list.Circles))
	require.Equal(t, 1, list.CurrentPage)
	require.Equal(t, 10, list.PageSize)
	require.NotEqual(t, 0, list.NumberOfPages)
	require.True(t, len(list.Circles) <= 10)

	for _, circle := range list.Circles {
		require.NotEmpty(t, circle.CircleID)
		require.NotEmpty(t, circle.Name)
		require.NotEmpty(t, circle.Description)
		require.NotEqual(t, 0, len(circle.Users))
		for _, user := range circle.Users {
			require.NotEmpty(t, user.UserID)
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.DisplayName)
		}
	}
}
