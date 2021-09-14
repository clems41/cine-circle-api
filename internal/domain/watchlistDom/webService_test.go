package watchlistDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/test"
	"cine-circle/pkg/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestHandler_AddMovie(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing user and movie to database
	userSample := sampler.GetUser()
	movie := sampler.GetMovie()

	correctBasePath := webServicePath + "/" + utils.IDToStr(movie.GetID())
	wrongBasePath := webServicePath + "/999999"

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodPost)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodPost)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodPost)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var watchlist repositoryModel.Watchlist
	err := DB.
		Preload("Movie").
		Preload("User").
		Take(&watchlist, "user_id = ? AND movie_id = ?", userSample.GetID(), movie.GetID()).
		Error
	require.NoError(t, err)
	require.Equal(t, userSample.GetID(), watchlist.User.GetID())
	require.Equal(t, movie.GetID(), watchlist.Movie.GetID())
}

func TestHandler_DeleteMovie(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing user and movie to database
	userSample := sampler.GetUser()
	movie := sampler.GetMovie()

	correctBasePath := webServicePath + "/" + utils.IDToStr(movie.GetID())
	wrongBasePath := webServicePath + "/999999"

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodDelete)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodDelete)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	var watchlist repositoryModel.Watchlist
	err := DB.
		Preload("Movie").
		Preload("User").
		Take(&watchlist, "user_id = ? AND movie_id = ?", userSample.GetID(), movie.GetID()).
		Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestHandler_AlreadyAdded(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing user to database with movie in its playlist
	userSample := sampler.GetUser()
	watchlist := sampler.GetWatchlist(userSample)

	// Get one element from user's watchlist
	element := test.RandomElement(watchlist)
	elem, ok := element.(repositoryModel.Watchlist)
	if !ok {
		t.Fatalf("Element should be type Watchlist")
	}

	correctBasePath := webServicePath + "/" + utils.IDToStr(elem.Movie.GetID())
	wrongBasePath := webServicePath + "/999999"

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var exists bool
	testingHTTPServer.DecodeResponse(resp, &exists)
	require.False(t, exists)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)
	testingHTTPServer.DecodeResponse(resp, &exists)
	require.True(t, exists)
}

func TestHandler_List(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)
	ruler := test.NewRuler(t)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing user and movie to database
	userSample := sampler.GetUser()

	// create watchlist for user with some movieSample
	watchlist := sampler.GetWatchlist(userSample)

	queryParameters := []test.KeyValue{
		{
			Key:   "page",
			Value: "1",
		},
		{
			Key:   "pageSize",
			Value: "15",
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with wrong movie id, should fail and return 404
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list List
	testingHTTPServer.DecodeResponse(resp, &list)
	require.Equal(t, list.NumberOfItems, len(watchlist))
	require.NotEqual(t, 0, len(list.Movies))
	require.Equal(t, list.CurrentPage, 1)
	require.NotEqual(t, list.NumberOfPages, 0)
	require.Equal(t, list.PageSize, 15)

	// TODO check number of movie

	for _, movie := range list.Movies {
		ruler.CheckStruct(movie, map[string]interface{}{
			"ID":          test.NotEmptyField{},
			"Title":       test.NotEmptyField{},
			"ReleaseDate": test.NotEmptyField{},
			"PosterPath":  test.NotEmptyField{},
		})
	}
}
