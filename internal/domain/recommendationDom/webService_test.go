package recommendationDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/test"
	"cine-circle/internal/utils"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Add existing users to database
	userNotInCircle1 := sampler.GetUser()
	userNotInCircle2 := sampler.GetUser()
	userInCircle := sampler.GetUser()

	// Creating circle with specific user
	circle := sampler.GetCircle(*userInCircle)

	movie := sampler.GetMovie()

	// fields for recommendation
	comment := fake.Sentences()
	fakeMovieId := uint(99999)
	fakeCircleId := uint(99999)
	fakeUserId := uint(99999)

	creation := Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   nil,
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userNotInCircle2)

	// Send request and check response code with wrong user authenticated, should fail and return 401
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Authenticate user for sending request (user not in circle), should fail and return 401
	testingHTTPServer.AuthenticateUserPermanently(userInCircle)

	// Send request and check response with wrong movieId, should fail and return 404
	creation = Creation{
		MovieID:   fakeMovieId,
		Comment:   comment,
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   []uint{userNotInCircle1.GetID()},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request and check response with empty comment, should fail and return 400
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   "",
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   []uint{userNotInCircle1.GetID()},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request and check response with empty recipient, should fail and return 400
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: nil,
		UserIDs:   nil,
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send request and check response with wrong userId, should fail and return 401 (like if users are not in same circle, event if user doesn't exists)
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: nil,
		UserIDs:   []uint{fakeUserId},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Send request and check response with wrong circleId, should fail and return 401 (like if user is not in circle, even if circle doesn't exists)
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: []uint{fakeCircleId},
		UserIDs:   nil,
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Send request and check response with user2 sending reco to user1 not in his contact list, should fail and return 401
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: nil,
		UserIDs:   []uint{userNotInCircle1.GetID()},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Creating circle with both user, should now be OK
	_ = sampler.GetCircle(*userNotInCircle1, *userNotInCircle2)

	// Send request and check response with user2 sending reco to user1 now in his contact list, should work
	testingHTTPServer.AuthenticateUserPermanently(userNotInCircle2)
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: nil,
		UserIDs:   []uint{userNotInCircle1.GetID()},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Send request and check response with user2 sending reco to other circle, should fail because user not in circle return 401
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   nil,
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Send request and check response with user2 sending reco to other circle, should work after adding him into circle
	err := DB.
		Exec("INSERT INTO circle_user (circle_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", circle.GetID(), userNotInCircle2.GetID()).
		Error
	require.NoError(t, err)
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   nil,
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Send request and check response with user2 sending reco to other circle and user1, should work
	err = DB.
		Exec("INSERT INTO circle_user (circle_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", circle.GetID(), userNotInCircle2.GetID()).
		Error
	require.NoError(t, err)
	creation = Creation{
		MovieID:   movie.GetID(),
		Comment:   comment,
		CircleIDs: []uint{circle.GetID()},
		UserIDs:   []uint{userNotInCircle1.GetID()},
	}
	resp = testingHTTPServer.SendRequestWithBody(webServicePath, http.MethodPost, creation)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check in DB if recommendations have been saved
	var recommendations []repositoryModel.Recommendation
	err = DB.
		Preload("Users").
		Preload("Movie").
		Preload("Sender").
		Preload("Circles").
		Preload("Circles.Users").
		Order("id").
		Find(&recommendations, "movie_id = ? AND sender_id = ?", movie.GetID(), userNotInCircle2.GetID()).
		Error
	require.NoError(t, err)
	require.Len(t, recommendations, 3, "should find 3 recommendations in DB")

	// Update circle with new users before checking recommendation details
	err = DB.Preload("Users").Take(&circle).Error
	require.NoError(t, err)

	for idx, reco := range recommendations {
		// Check details about movie
		require.True(t, reco.MovieID == reco.Movie.GetID())
		require.Equal(t, reco.Movie.Title, movie.Title)

		// Check details about sender
		require.True(t, reco.SenderID == reco.Sender.GetID())
		require.Equal(t, reco.Sender.Email, userNotInCircle2.Email)
		require.Equal(t, reco.Sender.Username, userNotInCircle2.Username)
		require.Equal(t, reco.Sender.DisplayName, userNotInCircle2.DisplayName)

		// Check details about circles and users
		switch idx {
		case 0:
			require.Len(t, reco.Users, 1)
			require.Len(t, reco.Circles, 0)
		case 1:
			require.Len(t, reco.Users, 0)
			require.Len(t, reco.Circles, 1)
			require.Len(t, reco.Circles[0].Users, len(circle.Users))
		case 2:
			require.Len(t, reco.Users, 1)
			require.Len(t, reco.Circles, 1)
			require.Len(t, reco.Circles[0].Users, len(circle.Users))
		}
	}
}

func TestHandler_List(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Create recommendations sent and received by user
	userSample := sampler.GetUser()
	sentRecommendations := sampler.GetRecommendationsSentByUser(userSample)
	receivedRecommendations := sampler.GetRecommendationsReceivedByUser(userSample)

	queryParameters := []test.KeyValue{
		{
			Key:   "field",
			Value: fake.Word(),
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	queryParameters = []test.KeyValue{
		{
			Key:   "field",
			Value: "date",
		},
		{
			Key:   "desc",
			Value: "true",
		},
	}

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view ViewList
	testingHTTPServer.DecodeResponse(resp, &view)
	require.Equal(t, len(receivedRecommendations)+len(sentRecommendations), view.NumberOfItems)
	require.Equal(t, len(view.Recommendations), view.NumberOfItems)
	require.Equal(t, view.CurrentPage, 1)
	require.NotEqual(t, view.NumberOfPages, 0)
	require.NotEqual(t, view.PageSize, 0)

	for _, recommendation := range view.Recommendations {
		require.NotEqual(t, 0, len(recommendation.Circles))
		for _, circle := range recommendation.Circles {
			require.NotEqual(t, 0, len(circle.Users))
			require.NotEqual(t, 0, circle.CircleID)
			require.NotEmpty(t, circle.Name)
			require.NotEmpty(t, circle.Description)
			for _, circleUser := range circle.Users {
				require.NotEmpty(t, circleUser.Username)
				require.NotEmpty(t, circleUser.DisplayName)
				require.NotEqual(t, 0, circleUser.UserID)
			}
		}
		require.NotEqual(t, 0, len(recommendation.Users))
		for _, user := range recommendation.Users {
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.DisplayName)
			require.NotEqual(t, 0, user.UserID)
		}
		require.NotEmpty(t, recommendation.Movie.Title)
		require.NotEmpty(t, recommendation.Movie.PosterPath)
		require.NotEqual(t, 0, recommendation.Movie.ID)
		require.NotEmpty(t, recommendation.Date)
		require.NotEmpty(t, recommendation.RecommendationType)
		require.NotEmpty(t, recommendation.Comment)
		require.NotEmpty(t, recommendation.Sender.Username)
		require.NotEmpty(t, recommendation.Sender.DisplayName)
		require.NotEqual(t, 0, recommendation.Sender.UserID)
	}
}

func TestHandler_List_ReceivedRecommendations(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Create recommendations sent and received by user
	userSample := sampler.GetUser()
	sampler.GetRecommendationsSentByUser(userSample)
	receivedRecommendations := sampler.GetRecommendationsReceivedByUser(userSample)

	queryParameters := []test.KeyValue{
		{
			Key:   "field",
			Value: fake.Word(),
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	queryParameters = []test.KeyValue{
		{
			Key:   "field",
			Value: "date",
		},
		{
			Key:   "desc",
			Value: "true",
		},
		{
			Key:   "recommendationType",
			Value: "received",
		},
	}

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view ViewList
	testingHTTPServer.DecodeResponse(resp, &view)
	require.Equal(t, len(receivedRecommendations), view.NumberOfItems)
	require.Equal(t, len(view.Recommendations), view.NumberOfItems)
	require.Equal(t, view.CurrentPage, 1)
	require.NotEqual(t, view.NumberOfPages, 0)
	require.NotEqual(t, view.PageSize, 0)

	for _, recommendation := range view.Recommendations {
		require.NotEqual(t, 0, len(recommendation.Circles))
		for _, circle := range recommendation.Circles {
			require.NotEqual(t, 0, len(circle.Users))
			require.NotEqual(t, 0, circle.CircleID)
			require.NotEmpty(t, circle.Name)
			require.NotEmpty(t, circle.Description)
			for _, circleUser := range circle.Users {
				require.NotEmpty(t, circleUser.Username)
				require.NotEmpty(t, circleUser.DisplayName)
				require.NotEqual(t, 0, circleUser.UserID)
			}
		}
		require.NotEqual(t, 0, len(recommendation.Users))
		for _, user := range recommendation.Users {
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.DisplayName)
			require.NotEqual(t, 0, user.UserID)
		}
		require.NotEmpty(t, recommendation.Movie.Title)
		require.NotEmpty(t, recommendation.Movie.PosterPath)
		require.NotEqual(t, 0, recommendation.Movie.ID)
		require.NotEmpty(t, recommendation.Date)
		require.NotEmpty(t, recommendation.RecommendationType)
		require.NotEmpty(t, recommendation.Comment)
		require.NotEmpty(t, recommendation.Sender.Username)
		require.NotEmpty(t, recommendation.Sender.DisplayName)
		require.NotEqual(t, 0, recommendation.Sender.UserID)
	}
}

func TestHandler_List_SentRecommendations(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Create recommendations sent and received by user
	userSample := sampler.GetUser()
	sentRecommendations := sampler.GetRecommendationsSentByUser(userSample)
	sampler.GetRecommendationsReceivedByUser(userSample)

	queryParameters := []test.KeyValue{
		{
			Key:   "field",
			Value: fake.Word(),
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	queryParameters = []test.KeyValue{
		{
			Key:   "field",
			Value: "date",
		},
		{
			Key:   "desc",
			Value: "true",
		},
		{
			Key:   "recommendationType",
			Value: "sent",
		},
	}

	// Send request and check response code with bad query params, should fail and return 400
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var view ViewList
	testingHTTPServer.DecodeResponse(resp, &view)
	require.Equal(t, len(sentRecommendations), view.NumberOfItems)
	require.Equal(t, len(view.Recommendations), view.NumberOfItems)
	require.Equal(t, view.CurrentPage, 1)
	require.NotEqual(t, view.NumberOfPages, 0)
	require.NotEqual(t, view.PageSize, 0)

	for _, recommendation := range view.Recommendations {
		require.NotEqual(t, 0, len(recommendation.Circles))
		for _, circle := range recommendation.Circles {
			require.NotEqual(t, 0, len(circle.Users))
			require.NotEqual(t, 0, circle.CircleID)
			require.NotEmpty(t, circle.Name)
			require.NotEmpty(t, circle.Description)
			for _, circleUser := range circle.Users {
				require.NotEmpty(t, circleUser.Username)
				require.NotEmpty(t, circleUser.DisplayName)
				require.NotEqual(t, 0, circleUser.UserID)
			}
		}
		require.NotEqual(t, 0, len(recommendation.Users))
		for _, user := range recommendation.Users {
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.DisplayName)
			require.NotEqual(t, 0, user.UserID)
		}
		require.NotEmpty(t, recommendation.Movie.Title)
		require.NotEmpty(t, recommendation.Movie.PosterPath)
		require.NotEqual(t, 0, recommendation.Movie.ID)
		require.NotEmpty(t, recommendation.Date)
		require.NotEmpty(t, recommendation.RecommendationType)
		require.NotEmpty(t, recommendation.Comment)
		require.NotEmpty(t, recommendation.Sender.Username)
		require.NotEmpty(t, recommendation.Sender.DisplayName)
		require.NotEqual(t, 0, recommendation.Sender.UserID)
	}
}

func TestHandler_List_SpecificMovie(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)

	// Create user and add it to some circles
	userSample := sampler.GetUser()
	movie := sampler.GetMovie()
	sampler.GetRecommendationsSentByUserForSpecificMovie(userSample, movie)
	receivedRecommendations := sampler.GetRecommendationsReceivedByUserForSpecificMovie(userSample, movie)

	queryParameters := []test.KeyValue{
		{
			Key:   "page",
			Value: "1",
		},
		{
			Key:   "pageSize",
			Value: "7",
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code, should fail with wrong movie id
	queryParameters = []test.KeyValue{
		{
			Key:   "page",
			Value: "1",
		},
		{
			Key:   "pageSize",
			Value: "7",
		},
		{
			Key:   "movieId",
			Value: "99999",
		},
	}
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Send request and check response code, should work

	queryParameters = []test.KeyValue{
		{
			Key:   "page",
			Value: "1",
		},
		{
			Key:   "pageSize",
			Value: "7",
		},
		{
			Key:   "movieId",
			Value: utils.IDToStr(movie.GetID()),
		},
	}
	resp = testingHTTPServer.SendRequestWithQueryParameters(webServicePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list ViewList
	testingHTTPServer.DecodeResponse(resp, &list)

	require.Equal(t, len(receivedRecommendations), list.NumberOfItems)
	require.NotEqual(t, 0, len(list.Recommendations))
	require.Equal(t, 1, list.CurrentPage)
	require.Equal(t, 7, list.PageSize)
	require.NotEqual(t, 0, list.NumberOfPages)
	require.True(t, len(list.Recommendations) <= 7)

	for _, recommendation := range list.Recommendations {
		require.NotEqual(t, 0, len(recommendation.Circles))
		for _, circle := range recommendation.Circles {
			require.NotEqual(t, 0, len(circle.Users))
			require.NotEqual(t, 0, circle.CircleID)
			require.NotEmpty(t, circle.Name)
			require.NotEmpty(t, circle.Description)
			for _, circleUser := range circle.Users {
				require.NotEmpty(t, circleUser.Username)
				require.NotEmpty(t, circleUser.DisplayName)
				require.NotEqual(t, 0, circleUser.UserID)
			}
		}
		require.NotEqual(t, 0, len(recommendation.Users))
		for _, user := range recommendation.Users {
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.DisplayName)
			require.NotEqual(t, 0, user.UserID)
		}
		require.Equal(t, movie.Title, recommendation.Movie.Title)
		require.Equal(t, movie.PosterPath, recommendation.Movie.PosterPath)
		require.Equal(t, movie.GetID(), recommendation.Movie.ID)
		require.NotEmpty(t, recommendation.Date)
		require.NotEmpty(t, recommendation.RecommendationType)
		require.NotEmpty(t, recommendation.Comment)
		require.NotEmpty(t, recommendation.Sender.Username)
		require.NotEmpty(t, recommendation.Sender.DisplayName)
		require.NotEqual(t, 0, recommendation.Sender.UserID)
	}
}

func TestHandler_ListUsers(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)
	basePath := webServicePath + "/users"

	// Create user and add it to some circles
	userSample := sampler.GetUser()
	var circles []*repositoryModel.Circle
	for range test.FakeRange(2, 6) {
		circles = append(circles, sampler.GetCircle(*userSample))
	}

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
	resp := testingHTTPServer.SendRequestWithQueryParameters(basePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code, should work
	resp = testingHTTPServer.SendRequestWithQueryParameters(basePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list UserList
	testingHTTPServer.DecodeResponse(resp, &list)
	require.NotEqual(t, 0, len(list.Users))
	require.Equal(t, list.CurrentPage, 1)
	require.NotEqual(t, list.NumberOfPages, 0)
	require.Equal(t, list.PageSize, 15)

	// calculate number of user in all circles
	var nbUsers int
	for _, circle := range circles {
		// we remove one to not count actual user
		nbUsers += len(circle.Users) - 1
	}
	require.Equal(t, list.NumberOfItems, nbUsers)
}

func TestHandler_ListUsers_Page2(t *testing.T) {
	DB, clean := test.OpenDatabase(t)
	defer clean()

	sampler := test.NewSampler(t, DB, true)

	webService := NewHandler(NewService(NewRepository(DB)))
	webServicePath := webService.WebServices()[0].RootPath()
	testingHTTPServer := test.NewTestingHTTPServer(t, webService)
	basePath := webServicePath + "/users"

	// Create user and add it to some circles
	userSample := sampler.GetUser()
	var circles []*repositoryModel.Circle
	for range test.FakeRange(2, 6) {
		circles = append(circles, sampler.GetCircle(*userSample))
	}

	queryParameters := []test.KeyValue{
		{
			Key:   "page",
			Value: "2",
		},
		{
			Key:   "pageSize",
			Value: "7",
		},
	}

	// Send request and check response code without authentication, should fail and return 401
	resp := testingHTTPServer.SendRequestWithQueryParameters(basePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	testingHTTPServer.AuthenticateUserPermanently(userSample)

	// Send request and check response code, should work
	resp = testingHTTPServer.SendRequestWithQueryParameters(basePath, http.MethodGet, queryParameters)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list UserList
	testingHTTPServer.DecodeResponse(resp, &list)
	require.NotEqual(t, 0, len(list.Users))
	require.Equal(t, list.CurrentPage, 2)
	require.NotEqual(t, list.NumberOfPages, 0)
	require.Equal(t, list.PageSize, 7)

	// calculate number of user in all circles
	var nbUsers int
	for _, circle := range circles {
		// we remove one to not count actual user
		nbUsers += len(circle.Users) - 1
	}
	require.Equal(t, list.NumberOfItems, nbUsers)
}
