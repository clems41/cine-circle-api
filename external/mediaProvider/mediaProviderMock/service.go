package mediaProviderMock

import (
	"cine-circle-api/external/mediaProvider"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"fmt"
	"github.com/icrowley/fake"
	"time"
)

var _ mediaProvider.Service = (*service)(nil)

type service struct {
}

func New() (svc *service) {
	return &service{}
}

func (svc *service) GetProviderName() (name string) {
	return "mediaProviderMock"
}

func (svc *service) Search(form mediaProvider.SearchForm) (medias []mediaProvider.MediaShort, total int64, err error) {
	total = int64(fakeData.FakeIntBetween(32, 100))
	for range fakeData.FakeRange(5, 32) {
		medias = append(medias, mediaProvider.MediaShort{
			Id:            fakeData.UuidWithOnlyAlphaNumeric(),
			Title:         fake.Title() + form.Keyword + fake.Title(),
			Language:      fake.Language(),
			OriginalTitle: fake.Title() + form.Keyword + fake.Title(),
			PosterUrl:     fake.StreetAddress(),
		})
	}
	return
}

func (svc *service) Get(mediaId string) (media mediaProvider.Media, err error) {
	if mediaId == "fake" { // Useful to test errors
		return media, fmt.Errorf("media %s cannot be found", media.Id)
	}
	var genres []string
	for range fakeData.FakeRange(1, 5) {
		genres = append(genres, fake.Word())
	}
	media = mediaProvider.Media{
		Id:            mediaId,
		Title:         fake.Title(),
		BackdropUrl:   fake.StreetAddress(),
		Genres:        fake.GetLangs(),
		Language:      fake.Language(),
		OriginalTitle: fake.Title(),
		Overview:      fake.Sentences(),
		PosterUrl:     fake.StreetAddress(),
		ReleaseDate:   fakeData.FakeTimeBefore(time.Now()),
		Runtime:       fakeData.FakeIntBetween(35, 236),
	}
	return
}
