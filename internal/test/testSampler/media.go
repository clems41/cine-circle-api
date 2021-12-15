package testSampler

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"time"
)

func (sampler *Sampler) GetCompletedMovie() (movie *model.Movie) {
	var genres []string
	for range fakeData.FakeRange(1, 5) {
		genres = append(genres, fake.Word())
	}
	movie = &model.Movie{
		Title:             fake.Title(),
		MediaProviderId:   fakeData.UuidWithOnlyAlphaNumeric(),
		MediaProviderName: "mediaProviderMock",
		Completed:         true,
		BackdropUrl:       fake.StreetAddress(),
		Genres:            genres,
		Language:          fake.Language(),
		OriginalTitle:     fake.Title(),
		Overview:          fake.Sentences(),
		PosterUrl:         fake.StreetAddress(),
		ReleaseDate:       fakeData.FakeTimeBefore(time.Now()),
		Runtime:           time.Duration(fakeData.FakeIntBetween(35, 236)) * time.Minute,
	}

	err := sampler.DB.Create(movie).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetUncompletedMovie() (movie *model.Movie) {
	var genres []string
	for range fakeData.FakeRange(1, 5) {
		genres = append(genres, fake.Word())
	}
	movie = &model.Movie{
		MediaProviderId:   fakeData.UuidWithOnlyAlphaNumeric(),
		MediaProviderName: "mediaProviderMock",
		Completed:         false,
	}

	err := sampler.DB.Create(movie).Error
	require.NoError(sampler.t, err)
	return
}
