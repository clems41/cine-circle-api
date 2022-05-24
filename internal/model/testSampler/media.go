package testSampler

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"time"
)

func (sampler *Sampler) GetMedia(customFields map[string]interface{}) (media *model.Media) {
	var genres []string
	for range fakeData.FakeRange(1, 5) {
		genres = append(genres, fake.Word())
	}
	media = &model.Media{
		Title:         fake.Title(),
		BackdropUrl:   fake.StreetAddress(),
		Genres:        genres,
		Language:      fake.Language(),
		OriginalTitle: fake.Title(),
		Overview:      fake.Sentences(),
		PosterUrl:     fake.StreetAddress(),
		ReleaseDate:   fakeData.FakeTimeBefore(time.Now()),
		Runtime:       fakeData.FakeIntBetween(35, 236),
		MediaType:     model.MovieMediaType,
	}
	err := fakeData.FillStructWithFields(media, customFields)
	require.NoError(sampler.t, err)
	return
}
