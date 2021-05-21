package movieDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/test"
	"cine-circle/internal/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler_Get(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	ruler := test.NewRuler(t)
	sampler := test.NewSampler(t, DB, false)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing circle to database
	userSample := sampler.GetUser()

	// ID of movie The Dark knight !!!
	movieId := uint(155)

	wrongBasePath := webServicePath + "/999999999"
	correctBasePath := webServicePath + "/" + utils.IDToStr(movieId)

	// Send request and check response code without authentication, should fail
	resp := testingHTTPServer.SendRequest(wrongBasePath, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (bad user)
	err := testingHTTPServer.AuthenticateUserPermanently(userSample)
	require.NoError(t, err, "User should be authenticated")

	// Send request with bas id, should return 404
	resp = testingHTTPServer.SendRequest(wrongBasePath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request with bas id, should return 302
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)

	var view View
	testingHTTPServer.DecodeResponse(resp, &view)

	// Check if view return correct movie
	ruler.CheckStruct(view, map[string]interface{}{
		"ID": movieId,
		"Title": test.NotEmptyField{},
		"ImdbId": test.NotEmptyField{},
		"BackdropPath": test.NotEmptyField{},
		"PosterPath": test.NotEmptyField{},
		"Genres": test.NotEmptyField{},
		"OriginalLanguage": test.NotEmptyField{},
		"OriginalTitle": "The Dark Knight",
		"Overview": test.NotEmptyField{},
		"ReleaseDate": test.NotEmptyField{},
		"Runtime": test.NotEmptyField{},
	})

	// Check if movie has been saved into DB
	var movie repositoryModel.Movie
	err = DB.
		Take(&movie, "id = ?", movieId).
		Error
	require.NoError(t, err)

	// Check if field's view are the same than movie in database
	ruler.CheckStruct(movie, map[string]interface{}{
		"ID": view.ID,
		"Title": view.Title,
		"ImdbId": view.ImdbId,
		"BackdropPath": view.BackdropPath,
		"PosterPath": view.PosterPath,
		"Genres": test.NotEmptyField{},
		"OriginalLanguage": view.OriginalLanguage,
		"OriginalTitle": view.OriginalTitle,
		"Overview": view.Overview,
		"ReleaseDate": view.ReleaseDate,
		"Runtime": view.Runtime,
	})
	require.Len(t, movie.Genres, len(view.Genres))
	for idx, movieGenre := range movie.Genres {
		require.Equal(t, movieGenre, view.Genres[idx])
	}

	// Check if movie can be retrieved a second time
	resp = testingHTTPServer.SendRequest(correctBasePath, http.MethodGet)
	require.Equal(t, http.StatusFound, resp.StatusCode)

	var view2 View
	testingHTTPServer.DecodeResponse(resp, &view2)

	ruler.CheckStruct(movie, map[string]interface{}{
		"ID": view.ID,
		"Title": view.Title,
		"ImdbId": view.ImdbId,
		"BackdropPath": view.BackdropPath,
		"PosterPath": view.PosterPath,
		"Genres": test.NotEmptyField{},
		"OriginalLanguage": view.OriginalLanguage,
		"OriginalTitle": view.OriginalTitle,
		"Overview": view.Overview,
		"ReleaseDate": view.ReleaseDate,
		"Runtime": view.Runtime,
	})
	require.Len(t, movie.Genres, len(view.Genres))
	for idx, movieGenre := range movie.Genres {
		require.Equal(t, movieGenre, view.Genres[idx])
	}
}
