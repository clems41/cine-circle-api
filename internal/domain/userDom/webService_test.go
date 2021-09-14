package userDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/test"
	utils2 "cine-circle/pkg/utils"
	"cine-circle/pkg/webService"
	"encoding/base64"
	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestHandler_CreateUser(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)

	userHandler := NewHandler(NewService(NewRepository(DB)))
	testingHTTPServer := test.NewTestingHTTPServer(t, userHandler)

	// Routes use for this test
	signUpBasePath := userHandler.WebServices()[0].RootPath() + "/sign-up"

	// Fields for creation
	username := fake.UserName()
	displayName := fake.FullName()
	password := test.FakePassword()
	email := fake.EmailAddress()

	// test creation with missing mandatory field : Username
	creation := Creation{
		Username: "",
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       email,
		},
		Password: password,
	}
	// Send request and check response code
	resp := testingHTTPServer.SendRequestWithBody(signUpBasePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test creation with missing mandatory field : DisplayName
	creation = Creation{
		Username: username,
		CommonFields: CommonFields{
			DisplayName: "",
			Email:       email,
		},
		Password: password,
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(signUpBasePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test creation with missing mandatory field : Password
	creation = Creation{
		Username: username,
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       email,
		},
		Password: "",
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(signUpBasePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test creation with missing mandatory field : Email
	creation = Creation{
		Username: username,
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       "",
		},
		Password: password,
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(signUpBasePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test creation with all correct fields
	creation = Creation{
		Username: username,
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       email,
		},
		Password: password,
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(signUpBasePath, http.MethodPost, creation)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"UserID":      test.NotEmptyField{},
		"DisplayName": creation.DisplayName,
		"Username":    strings.ToLower(creation.Username),
	})

	// Check that are not returned with View object
	var user repositoryModel.User
	err := DB.
		Find(&user, "id = ?", view.UserID).
		Error
	require.NoError(t, err, "Should not return error but got %v", err)
	require.Equal(t, creation.Email, user.Email)

	// check if password has been correctly salt and hash
	assert.NotEqual(t, creation.Password, user.HashedPassword, "password should be hashed")
	err = utils2.CompareHashAndPassword(user.HashedPassword, creation.Password)
	require.NoError(t, err, "passwords should be the same (using hash comparison) but got %v", err)
}

func TestHandler_Update(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// Add existing user to database
	userSample := sampler.GetUser()

	// Routes use for this test
	usersBasePath := userWebService.WebServices()[1].RootPath()

	// New fields to update user with
	displayName := fake.FullName()
	email := fake.EmailAddress()

	// test update without authentication
	update := Update{
		CommonFields: CommonFields{
			DisplayName: "",
			Email:       "",
		},
	}
	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequestWithBody(usersBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// test update with missing mandatory field : DisplayName and Email
	update = Update{
		CommonFields: CommonFields{
			DisplayName: "",
			Email:       "",
		},
	}
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(usersBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test update with only email
	update = Update{
		CommonFields: CommonFields{
			DisplayName: "",
			Email:       email,
		},
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(usersBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test update with only displayName
	update = Update{
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       "",
		},
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(usersBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test update with all fields
	update = Update{
		CommonFields: CommonFields{
			DisplayName: displayName,
			Email:       email,
		},
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(usersBasePath, http.MethodPut, update)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"UserID":      test.NotEmptyField{},
		"DisplayName": update.DisplayName,
		"Username":    strings.ToLower(*userSample.Username),
	})

	// Check that are not returned with View object
	var user repositoryModel.User
	err := DB.
		Find(&user, "id = ?", view.UserID).
		Error
	require.NoError(t, err, "Should not return error but got %v", err)
	require.Equal(t, update.Email, user.Email)
}

func TestHandler_UpdatePassword(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// New fields to update user with
	oldPassword := test.FakePassword()
	newPassword := test.FakePassword()

	// Add existing user to database
	userSample := sampler.GetUserWithSpecificPassword(oldPassword)

	// Define testing base path
	testingBasePath := webServicePath + "/password"

	// test updatePassword without authentication
	updatePassword := UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequestWithBody(testingBasePath, http.MethodPut, updatePassword)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// test updatePassword with missing mandatory field : OldPassword
	updatePassword = UpdatePassword{
		OldPassword: "",
		NewPassword: newPassword,
	}
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(testingBasePath, http.MethodPut, updatePassword)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test updatePassword with wrong mandatory field : OldPassword
	updatePassword = UpdatePassword{
		OldPassword: "wrongPassword",
		NewPassword: newPassword,
	}
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(testingBasePath, http.MethodPut, updatePassword)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// test updatePassword with missing mandatory field : NewPassword
	updatePassword = UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: "",
	}
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(testingBasePath, http.MethodPut, updatePassword)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// test updatePassword with all correct fields
	updatePassword = UpdatePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithBody(testingBasePath, http.MethodPut, updatePassword)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// check if password has been correctly salt and hash
	var user repositoryModel.User
	err := DB.
		Select("hashed_password").
		Find(&user, "id = ?", view.UserID).
		Error
	require.NoError(t, err, "Should not return error but got %v", err)
	assert.NotEqual(t, updatePassword.NewPassword, user.HashedPassword, "password should be hashed")
	assert.NotEqual(t, updatePassword.OldPassword, user.HashedPassword, "password should be hashed")
	err = utils2.CompareHashAndPassword(user.HashedPassword, updatePassword.OldPassword)
	require.Error(t, err, "passwords should not be the same (oldPassword)")
	err = utils2.CompareHashAndPassword(user.HashedPassword, updatePassword.NewPassword)
	require.NoError(t, err, "passwords should be the same (using hash comparison) but got %v", err)
}

func TestHandler_Delete(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// Add existing user to database
	userSample := sampler.GetUser()

	// Check if user has been correctly created
	var user repositoryModel.User
	err := DB.
		Take(&user, "id = ?", userSample.GetID()).
		Error
	require.NoError(t, err)
	require.Equal(t, *userSample.Username, *user.Username)
	require.Equal(t, userSample.Email, user.Email)
	require.Equal(t, userSample.HashedPassword, user.HashedPassword)
	require.Equal(t, userSample.DisplayName, user.DisplayName)

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(webServicePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code
	resp = testingHTTPServer.SendRequest(webServicePath, http.MethodDelete)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Check if user has been correctly deleted
	err = DB.
		Take(&user, "id = ?", userSample.GetID()).
		Error
	require.Error(t, err)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestHandler_Get(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// Add existing user to database
	userSample := sampler.GetUser()

	// Create different testing base path
	wrongID := test.FakeIntBetween(9999, 99999999)
	testingBasePath := webServicePath + "/" + utils2.IDToStr(userSample.GetID())
	wrongTestingBasePath := webServicePath + "/" + strconv.Itoa(wrongID)

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(wrongTestingBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code
	resp = testingHTTPServer.SendRequest(wrongTestingBasePath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request and check response code
	resp = testingHTTPServer.SendRequest(testingBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)
	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"UserID":      userSample.GetID(),
		"DisplayName": userSample.DisplayName,
		"Username":    *userSample.Username,
	})
}

func TestHandler_SearchUsers(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// Variables used for searching on username
	commonUsernamePart := (fake.UserName() + fake.UserName())[:3]
	matchingUsername1 := fake.UserName() + commonUsernamePart + fake.UserName()
	matchingUsername2 := commonUsernamePart + fake.UserName()
	matchingUsername3 := fake.UserName() + commonUsernamePart

	// Variables used for searching on displayName
	commonDisplayNamePart := fake.FullName()[3:9]
	matchingDisplayName1 := fake.FullName() + commonDisplayNamePart + fake.FullName()
	matchingDisplayName2 := commonDisplayNamePart + fake.FullName()
	matchingDisplayName3 := fake.FullName() + commonDisplayNamePart

	// Create some users using matching and not matching fields
	fakeUsername1 := fake.UserName()
	fakeUsername2 := fake.UserName()
	users := []repositoryModel.User{
		{
			Username:       &matchingUsername2,
			DisplayName:    matchingDisplayName3,
			Email:          fake.EmailAddress(),
			HashedPassword: test.FakePassword(),
		},
		{
			Username:       &matchingUsername1,
			DisplayName:    fake.FullName(),
			Email:          fake.EmailAddress(),
			HashedPassword: test.FakePassword(),
		},
		{
			Username:       &fakeUsername1,
			DisplayName:    matchingDisplayName1,
			Email:          fake.EmailAddress(),
			HashedPassword: test.FakePassword(),
		},
		{
			Username:       &matchingUsername3,
			DisplayName:    matchingDisplayName2,
			Email:          fake.EmailAddress(),
			HashedPassword: test.FakePassword(),
		},
		{
			Username:       &fakeUsername2,
			DisplayName:    fake.FullName(),
			Email:          fake.EmailAddress(),
			HashedPassword: test.FakePassword(),
		},
	}
	// Save users into database
	err := DB.
		Create(&users).
		Error
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Send request and check response code, should fail because user is not authenticate
	queryParameters := []test.KeyValue{
		{
			Key:   "search",
			Value: "a",
		},
	}
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	userSample := sampler.GetUser()
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Check with keyword matching nothing : should return not error and length of 0
	queryParameters[0].Value = fake.UserName()
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var views []View
	testingHTTPServer.DecodeResponse(resp, &views)
	require.Len(t, views, 0)

	// Check that search on email is not working : should return not error and length of 0
	queryParameters[0].Value = users[2].Email
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	testingHTTPServer.DecodeResponse(resp, &views)
	require.Len(t, views, 0)

	// Check with keyword matching username : should return not error and length of 3
	queryParameters[0].Value = commonUsernamePart
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	testingHTTPServer.DecodeResponse(resp, &views)
	require.Len(t, views, 3)

	// Check with keyword matching displayName : should return not error and length of 3
	queryParameters[0].Value = commonDisplayNamePart
	// Send request and check response code
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	testingHTTPServer.DecodeResponse(resp, &views)
	require.Len(t, views, 3)
}

func TestHandler_GenerateToken(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)
	testingBasePath := webServicePath + "/sign-in"

	// Create userSample used for getting token
	password := test.FakePassword()
	userSample := sampler.GetUserWithSpecificPassword(password)

	// Create correct and wrong values for header
	fakeUsername := strings.ToLower(fake.UserName())
	fakePassword := test.FakePassword()
	fakeUsernameWithCorrectPasswordEncoded := base64.StdEncoding.EncodeToString([]byte(
		fakeUsername + constant.UsernamePasswordDelimiterForHeader + password))
	fakeUsernameWithFakePasswordEncoded := base64.StdEncoding.EncodeToString([]byte(
		fakeUsername + constant.UsernamePasswordDelimiterForHeader + fakePassword))
	correctUsernameWithFakePasswordEncoded := base64.StdEncoding.EncodeToString([]byte(
		*userSample.Username + constant.UsernamePasswordDelimiterForHeader + fakePassword))
	correctUsernameWithCorrectPasswordEncoded := base64.StdEncoding.EncodeToString([]byte(
		*userSample.Username + constant.UsernamePasswordDelimiterForHeader + password))

	// Try to getting token with wrong username and wrong password
	headerParameters := []test.KeyValue{
		{
			Key:   constant.AuthenticationHeaderName,
			Value: constant.AuthenticationHeaderPrefixValue + fakeUsernameWithFakePasswordEncoded,
		},
	}
	resp := testingHTTPServer.SendRequestWithHeaders(testingBasePath, http.MethodPost, headerParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try to getting token with wrong username and correct password
	headerParameters[0].Value = constant.AuthenticationHeaderPrefixValue + fakeUsernameWithCorrectPasswordEncoded
	resp = testingHTTPServer.SendRequestWithHeaders(testingBasePath, http.MethodPost, headerParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try to getting token with correct username and wrong password
	headerParameters[0].Value = constant.AuthenticationHeaderPrefixValue + correctUsernameWithFakePasswordEncoded
	resp = testingHTTPServer.SendRequestWithHeaders(testingBasePath, http.MethodPost, headerParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try to getting token with correct username and correct password
	headerParameters[0].Value = constant.AuthenticationHeaderPrefixValue + correctUsernameWithCorrectPasswordEncoded
	resp = testingHTTPServer.SendRequestWithHeaders(testingBasePath, http.MethodPost, headerParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var token string
	testingHTTPServer.DecodeResponse(resp, &token)
	_, err := webService.CheckToken(token)
	require.NoError(t, err, "Token should be valid")
}

func TestHandler_GetOwnUserInfo(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)
	ruler := test.NewRuler(t)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)
	testingBasePath := webServicePath + "/me"

	// Send request and check response code, should fail because user is not authenticate
	resp := testingHTTPServer.SendRequest(testingBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request
	userSample := sampler.GetUser()
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code
	resp = testingHTTPServer.SendRequest(testingBasePath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var view ViewMe
	testingHTTPServer.DecodeResponse(resp, &view)
	// Check struct returned
	ruler.CheckStruct(view, map[string]interface{}{
		"UserID":      userSample.GetID(),
		"DisplayName": userSample.DisplayName,
		"Username":    *userSample.Username,
		"Email":       userSample.Email,
	})
}

func TestHandler_UsernameExists(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, false)

	userWebService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := userWebService.WebServices()[1].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, userWebService)

	// Add existing user to database
	userSample := sampler.GetUser()

	// Define route test
	fakeUsername := fake.UserName()
	wrongExistsBasePath := webServicePath + "/" + fakeUsername + "/exists"
	existsBasePath := webServicePath + "/" + *userSample.Username + "/exists"

	// Send request and check response code for wrong username
	resp := testingHTTPServer.SendRequest(wrongExistsBasePath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	var exists bool
	testingHTTPServer.DecodeResponse(resp, &exists)
	require.False(t, exists, "This username should not exists")

	// Send request and check response code for correct username
	resp = testingHTTPServer.SendRequest(existsBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)
	testingHTTPServer.DecodeResponse(resp, &exists)
	require.True(t, exists, "This username should exists")
}
