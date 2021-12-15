package mediaDom

import (
	"cine-circle-api/internal/repository/instance/mediaRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/internal/service/mediaProvider/mediaProviderMock"
	"cine-circle-api/internal/test/setupTestCase"
	"cine-circle-api/internal/test/testSampler"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestHandler_Search(t *testing.T) {
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	testPath := basePath
	user := sampler.GetUser()
	keyword := fake.Words()
	queryParameters := map[string][]string{
		keywordQueryParameter.Name: {keyword},
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with authenticated user --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that view is not empty
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"Page": map[string]interface{}{
			"NumberOfItems": testRuler.NotEmptyField{},
			"NumberOfPages": testRuler.NotEmptyField{},
			"PageSize":      testRuler.NotEmptyField{},
			"CurrentPage":   testRuler.NotEmptyField{},
		},
		"Result": testRuler.NotEmptyField{},
	})
	require.True(t, len(view.Result) > 0)

	// Check that all movies in result have been stored in database and marked as no completed (because not all fields have been filled)
	var ids []uint
	for _, result := range view.Result {
		ids = append(ids, result.Id)
	}
	var movies []model.Movie
	err := db.Find(&movies, "id in ?", ids).Error
	require.NoError(t, err)

	require.Equal(t, len(movies), len(view.Result))

	for _, movie := range movies {
		require.Equal(t, "mediaProviderMock", movie.MediaProviderName)
		require.NotEqual(t, "", movie.MediaProviderId)
		require.Equal(t, false, movie.Completed)
	}
}

// TestHandler_SearchThenGet will check the following workflow
//  - search for movies
//  - save all result into database with provider name and id
//  - mark result as uncompleted
//  - try to get one movie from result
//  - get movie from provider using previous provider id stored in database
//  - store now movie with all info from provider into database
func TestHandler_SearchThenGet(t *testing.T) {
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	testPath := basePath
	user := sampler.GetUser()
	keyword := fake.Words()
	queryParameters := map[string][]string{
		keywordQueryParameter.Name: {keyword},
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with authenticated user --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that view is not empty
	var searchView SearchView
	httpMock.DecodeResponse(resp, &searchView)
	require.True(t, len(searchView.Result) > 0)

	// Now try to get this movie
	uncompletedExistingMovieId := searchView.Result[0].Id
	getPath := fmt.Sprintf("%s/%d", basePath, uncompletedExistingMovieId)

	resp = httpMock.SendRequest(getPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that now movie is completed
	var movie model.Movie
	err := db.Take(&movie, uncompletedExistingMovieId).Error
	require.NoError(t, err)
	require.Equal(t, true, movie.Completed)

	// Check view fields
	var view GetView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"Id":            movie.ID,
		"Title":         movie.Title,
		"BackdropUrl":   movie.BackdropUrl,
		"Genres":        testRuler.NotEmptyField{},
		"Language":      movie.Language,
		"OriginalTitle": movie.OriginalTitle,
		"Overview":      movie.Overview,
		"PosterUrl":     movie.PosterUrl,
		"ReleaseDate":   movie.ReleaseDate,
		"Runtime":       movie.Runtime,
	})
}

// TestHandler_Get_AlreadyCompletedMovie will test that movie returned is the one already stored in database
func TestHandler_Get_AlreadyCompletedMovie(t *testing.T) {
	_, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	existingMovie := sampler.GetCompletedMovie()
	wrongMovieId := 999
	correctPath := fmt.Sprintf("%s/%d", basePath, existingMovie.ID)
	wrongPath := fmt.Sprintf("%s/%d", basePath, wrongMovieId)
	user := sampler.GetUser()

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with authenticated user but non-existing movie --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongPath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with authenticated user --> OK 200
	resp = httpMock.SendRequest(correctPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that view is not empty
	var view GetView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"Id":            existingMovie.ID,
		"Title":         existingMovie.Title,
		"BackdropUrl":   existingMovie.BackdropUrl,
		"Genres":        testRuler.NotEmptyField{},
		"Language":      existingMovie.Language,
		"OriginalTitle": existingMovie.OriginalTitle,
		"Overview":      existingMovie.Overview,
		"PosterUrl":     existingMovie.PosterUrl,
		"ReleaseDate":   existingMovie.ReleaseDate,
		"Runtime":       existingMovie.Runtime,
	})
}

// TestHandler_Get_UncompletedMovie will test that movie info are filled when uncompleted movie is requested
func TestHandler_Get_UncompletedMovie(t *testing.T) {
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	existingMovie := sampler.GetUncompletedMovie()
	wrongMovieId := 999
	correctPath := fmt.Sprintf("%s/%d", basePath, existingMovie.ID)
	wrongPath := fmt.Sprintf("%s/%d", basePath, wrongMovieId)
	user := sampler.GetUser()

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(correctPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with authenticated user but non-existing movie --> NOK 404
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(wrongPath, http.MethodGet)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with authenticated user --> OK 200
	resp = httpMock.SendRequest(correctPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that now movie is completed
	err := db.Take(&existingMovie).Error
	require.NoError(t, err)
	require.Equal(t, true, existingMovie.Completed)

	// Check view fields
	var view GetView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"Id":            existingMovie.ID,
		"Title":         existingMovie.Title,
		"BackdropUrl":   existingMovie.BackdropUrl,
		"Genres":        testRuler.NotEmptyField{},
		"Language":      existingMovie.Language,
		"OriginalTitle": existingMovie.OriginalTitle,
		"Overview":      existingMovie.Overview,
		"PosterUrl":     existingMovie.PosterUrl,
		"ReleaseDate":   existingMovie.ReleaseDate,
		"Runtime":       existingMovie.Runtime,
	})
}

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = setupTestCase.OpenCleanDatabaseFromTemplate(t)
	repo := mediaRepository.New(db)
	mock := mediaProviderMock.New()
	svc := NewService(mock, repo)
	ws := NewHandler(svc)
	httpMock = httpServerMock.New(t, logger.Logger(), ws)
	sampler = testSampler.New(t, db, populateDatabase)
	ruler = testRuler.New(t)
	return
}
