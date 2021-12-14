package mediaDom

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/internal/service/mediaProvider"
)

var _ Service = (*service)(nil)

type Service interface {
	Get(form GetForm) (view GetView, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type Repository interface {
	Get(mediaId uint) (media model.Media, ok bool, err error)
	Save(media *model.Media) (err error)
}

type service struct {
	mediaProvider mediaProvider.Service
	repository    Repository
}

func NewService(mediaProvider mediaProvider.Service, repository Repository) Service {
	return &service{
		repository:    repository,
		mediaProvider: mediaProvider,
	}
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	return
}
