package mediaProviderMock

import (
	mediaProvider2 "cine-circle-api/external/mediaProvider"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"fmt"
	"github.com/icrowley/fake"
	"time"
)

var _ mediaProvider2.Service = (*service)(nil)

type service struct {
}

func New() (svc *service) {
	return &service{}
}

func (svc *service) GetProviderName() (name string) {
	return "mediaProviderMock"
}

func (svc *service) Search(form mediaProvider2.SearchForm) (view mediaProvider2.SearchView, err error) {
	view.CurrentPage = form.Page
	view.NumberOfPages = fakeData.FakeIntBetween(1, 10)
	view.NumberOfItems = fakeData.FakeIntBetween(1, 100)
	for range fakeData.FakeRange(5, 32) {
		view.Result = append(view.Result, mediaProvider2.MovieShortView{
			Id:            fakeData.UuidWithOnlyAlphaNumeric(),
			Title:         fake.Title() + form.Keyword + fake.Title(),
			Language:      fake.Language(),
			OriginalTitle: fake.Title() + form.Keyword + fake.Title(),
			PosterUrl:     fake.StreetAddress(),
		})
	}
	return
}

func (svc *service) Get(form mediaProvider2.MovieForm) (view mediaProvider2.MovieView, err error) {
	if form.Id == "fake" { // Useful to test errors
		return view, fmt.Errorf("movie %s cannot be found", form.Id)
	}
	var genres []string
	for range fakeData.FakeRange(1, 5) {
		genres = append(genres, fake.Word())
	}
	view = mediaProvider2.MovieView{
		Id:            form.Id,
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
