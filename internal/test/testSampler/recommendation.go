package testSampler

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"time"
)

func (sampler *Sampler) GetRecommendation() (recommendation *model.Recommendation) {
	movie := sampler.GetCompletedMovie()
	sender := sampler.GetUser()
	recommendation = &model.Recommendation{
		SenderId: sender.ID,
		Circles:  nil,
		MovieId:  movie.ID,
		Text:     fake.Sentences(),
		Date:     time.Now(),
	}

	for range fakeData.FakeRange(1, 3) {
		circle := sampler.GetCircle()
		require.NotNil(sampler.t, circle)
		recommendation.Circles = append(recommendation.Circles, *circle)
	}

	err := sampler.DB.Create(recommendation).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetRecommendationReceivedBySpecificCircle(circle *model.Circle) (recommendation *model.Recommendation) {
	movie := sampler.GetCompletedMovie()
	sender := sampler.GetUser()
	recommendation = &model.Recommendation{
		SenderId: sender.ID,
		Circles:  nil,
		MovieId:  movie.ID,
		Text:     fake.Sentences(),
		Date:     time.Now(),
	}
	require.NotNil(sampler.t, circle)
	recommendation.Circles = append(recommendation.Circles, *circle)

	for range fakeData.FakeRange(1, 3) {
		fakeCircle := sampler.GetCircle()
		require.NotNil(sampler.t, fakeCircle)
		recommendation.Circles = append(recommendation.Circles, *fakeCircle)
	}

	err := sampler.DB.Create(recommendation).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetRecommendationSentBySpecificUser(sender *model.User) (recommendation *model.Recommendation) {
	movie := sampler.GetCompletedMovie()
	recommendation = &model.Recommendation{
		SenderId: sender.ID,
		Circles:  nil,
		MovieId:  movie.ID,
		Text:     fake.Sentences(),
		Date:     time.Now(),
	}

	for range fakeData.FakeRange(1, 3) {
		circle := sampler.GetCircle()
		require.NotNil(sampler.t, circle)
		recommendation.Circles = append(recommendation.Circles, *circle)
	}

	err := sampler.DB.Create(recommendation).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetRecommendationSentByUserWithSpecificMovie(sender *model.User, movie *model.Movie) (recommendation *model.Recommendation) {
	recommendation = &model.Recommendation{
		SenderId: sender.ID,
		Circles:  nil,
		MovieId:  movie.ID,
		Text:     fake.Sentences(),
		Date:     time.Now(),
	}

	for range fakeData.FakeRange(1, 3) {
		circle := sampler.GetCircle()
		require.NotNil(sampler.t, circle)
		recommendation.Circles = append(recommendation.Circles, *circle)
	}

	err := sampler.DB.Create(recommendation).Error
	require.NoError(sampler.t, err)
	return
}
