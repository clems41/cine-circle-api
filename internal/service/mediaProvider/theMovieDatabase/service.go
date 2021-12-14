package theMovieDatabase

import "cine-circle-api/internal/service/mediaProvider"

var _ mediaProvider.Service = (*service)(nil)

type service struct {
}

func New() (svc *service) {
	return &service{}
}

func (svc *service) Search(form mediaProvider.SearchForm) (view mediaProvider.SearchView, err error) {
	return
}

func (svc *service) Get(form mediaProvider.MediaForm) (view mediaProvider.MediaView, err error) {
	return
}
