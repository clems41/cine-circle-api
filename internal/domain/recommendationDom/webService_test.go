package recommendationDom

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/model/testSampler"
	"cine-circle-api/internal/repository"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

/* Create */

func TestHandler_Send(t *testing.T) {
	testPath := basePath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	sender := sampler.GetUser()
	movie := sampler.GetCompletedMovie()
	var circles []*model.Circle
	var circleIds []uint
	for range fakeData.FakeRange(1, 3) {
		circle := sampler.GetCircleWithUsers()
		circles = append(circles, circle)
		circleIds = append(circleIds, circle.ID)
	}
	correctForm := SendForm{CommonForm{
		CirclesIds: circleIds,
		MediaId:    movie.ID,
		Text:       fake.Sentences(),
	}}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with wrong form (missing text) --> NOK 400
	wrongForm := correctForm
	wrongForm.Text = ""
	httpMock.AuthenticateUserPermanently(sender)
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with wrong form (missing movieId) --> NOK 400
	wrongForm = correctForm
	wrongForm.MediaId = 0
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with wrong form (non-existing movieId) --> NOK 404
	wrongForm = correctForm
	wrongForm.MediaId = 999
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrong form (non-existing circleId) --> NOK 404
	wrongForm = correctForm
	wrongForm.CirclesIds = []uint{999}
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Try with wrong form (missing CirclesIds) --> NOK 400
	wrongForm = correctForm
	wrongForm.CirclesIds = nil
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct form --> OK 201
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check view fields
	var view SendView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id": testRuler.NotEmptyField{},
			"Sender": map[string]interface{}{
				"Id":        sender.ID,
				"Firstname": sender.FirstName,
				"Lastname":  sender.LastName,
				"Username":  sender.Username,
			},
			"Circles": testRuler.NotEmptyField{},
			"Movie": map[string]interface{}{
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
			},
			"Text": correctForm.Text,
			"Date": testRuler.NotEmptyField{},
			"Type": recommendationConst.SentType,
		},
	})
	require.Len(t, view.Circles, len(circles))
	require.True(t, view.Date.Before(time.Now()))
	for _, circleView := range view.Circles {
		require.NotEqual(t, 0, len(circleView.Users))
	}

	// Check that circle has been created into database
	var recommendation model.Recommendation
	err := db.
		Preload("Movie").
		Preload("Circles").
		Preload("Circles.Users").
		Preload("Sender").
		Take(&recommendation, view.Id).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, recommendation.ID)
	require.Len(t, recommendation.Circles, len(circles))
	require.Equal(t, recommendation.Sender.ID, sender.ID)
	require.Equal(t, recommendation.Movie.ID, movie.ID)
	for _, circleDb := range recommendation.Circles {
		require.NotEqual(t, 0, len(circleDb.Users))
	}
}

/* Search */

// TestHandler_Search check that circles returned are only one with authenticated user in it
func TestHandler_Search(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	movie := sampler.GetCompletedMovie()
	nbSentRecommendations := 6
	nbRecommendationsForMovie := 2
	nbReceivedRecommendations := 4
	var recommendations []*model.Recommendation
	for idx := range fakeData.FakeRange(14, 20) { // create between 14 and 20 recommendations (random number) but only 6 + 4 + 2 will match with user
		if idx < nbSentRecommendations {
			recommendations = append(recommendations, sampler.GetRecommendationSentBySpecificUser(user))
		} else if idx < nbSentRecommendations+nbRecommendationsForMovie {
			recommendations = append(recommendations, sampler.GetRecommendationSentByUserWithSpecificMovie(user, movie))
		} else if idx < nbSentRecommendations+nbRecommendationsForMovie+nbReceivedRecommendations {
			circle := sampler.GetCircleWithSpecificUser(user)
			recommendations = append(recommendations, sampler.GetRecommendationReceivedBySpecificCircle(circle))
		} else {
			recommendations = append(recommendations, sampler.GetRecommendation())
		}
	}
	page := 1
	pageSize := 15

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, recommendationConst.AllType, 0))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user with wrong type --> NOK 400
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, fake.Words(), 0))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with user with no filters
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, recommendationConst.AllType, 0))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbReceivedRecommendations+nbSentRecommendations+nbRecommendationsForMovie, view.NumberOfItems)
	require.Equal(t, nbReceivedRecommendations+nbSentRecommendations+nbRecommendationsForMovie, len(view.Recommendations))

	// Try with user with sent filter
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, recommendationConst.SentType, 0))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbSentRecommendations+nbRecommendationsForMovie, view.NumberOfItems)
	require.Equal(t, nbSentRecommendations+nbRecommendationsForMovie, len(view.Recommendations))
	for _, recommendation := range view.Recommendations {
		require.Equal(t, user.ID, recommendation.Sender.Id)
	}

	// Try with user with received filter
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, recommendationConst.ReceivedType, 0))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbReceivedRecommendations, view.NumberOfItems)
	require.Equal(t, nbReceivedRecommendations, len(view.Recommendations))
	for _, recommendation := range view.Recommendations {
		var userFound bool
		for _, circle := range recommendation.Circles {
			for _, userCircle := range circle.Users {
				if userCircle.Id == user.ID {
					userFound = true
				}
			}
		}
		require.True(t, userFound, "user %d should be in one of circle from recommendation", user.ID)
	}

	// Try with user with movie
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, recommendationConst.AllType, movie.ID))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbRecommendationsForMovie, view.NumberOfItems)
	require.Equal(t, nbRecommendationsForMovie, len(view.Recommendations))
	for _, recommendation := range view.Recommendations {
		require.Equal(t, movie.ID, recommendation.Movie.Id)
	}
}

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = setupTestCase.OpenCleanDatabaseFromTemplate(t)
	repo := repository.New(db)
	userRepo := repository.New(db)
	mediaRepo := repository.New(db)
	circleRepo := repository.New(db)
	svc := NewService(repo, userRepo, mediaRepo, circleRepo)
	ws := NewHandler(svc)
	httpMock = httpServerMock.New(t, logger.Logger(), ws)
	sampler = testSampler.New(t, db, populateDatabase)
	ruler = testRuler.New(t)
	return
}

// searchQueryParameters will return queryParameters based on search parameters
func searchQueryParameters(page, pageSize int, recommendationType string, mediaId uint) (queryParameters map[string][]string) {
	queryParameters = map[string][]string{
		pageQueryParameter.Name:     {fmt.Sprintf("%d", page)},
		pageSizeQueryParameter.Name: {fmt.Sprintf("%d", pageSize)},
	}
	if recommendationType != "" {
		queryParameters[recommendationTypeQueryParameter.Name] = []string{recommendationType}
	}
	if mediaId != 0 {
		queryParameters[mediaIdQueryParameter.Name] = []string{fmt.Sprintf("%d", mediaId)}
	}
	return
}
