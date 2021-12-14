package exempleDom

import (
	"cine-circle-api/internal/repository/instance/exempleRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(form CreateForm) (view CreateView, err error)
	Update(form UpdateForm) (view UpdateView, err error)
	Get(form GetForm) (view GetView, err error)
	Delete(form DeleteForm) (err error)
	Search(form SearchForm) (view SearchView, err error)
}

type Repository interface {
	Get(exempleId uint) (exemple model.Exemple, ok bool, err error)
	Create(exemple *model.Exemple) (err error)
	Save(exemple *model.Exemple) (err error)
	Delete(exempleId uint) (err error)
	Search(repoForm exempleRepository.SearchForm) (repoView exempleRepository.SearchView, err error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (svc *service) Create(form CreateForm) (view CreateView, err error) {
	exemple := svc.fromFormToModel(form.CommonForm)
	err = svc.repository.Create(&exemple)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(exemple)
	return
}

func (svc *service) Update(form UpdateForm) (view UpdateView, err error) {
	exemple, ok, err := svc.repository.Get(form.ExempleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errExempleIntrouvable
	}
	exemple = svc.fromFormToModel(form.CommonForm)
	exemple.ID = form.ExempleId
	err = svc.repository.Save(&exemple)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(exemple)
	return
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	exemple, ok, err := svc.repository.Get(form.ExempleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errExempleIntrouvable
	}
	view.CommonView = svc.fromModelToView(exemple)
	return
}

func (svc *service) Delete(form DeleteForm) (err error) {
	_, ok, err := svc.repository.Get(form.ExempleId)
	if err != nil {
		return
	}
	if !ok {
		return errExempleIntrouvable
	}

	err = svc.repository.Delete(form.ExempleId)
	if err != nil {
		return
	}
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	repoForm := exempleRepository.SearchForm{
		PaginationQuery: gormUtils.PaginationQuery{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		SortQuery: gormUtils.FromSortRequestToSortQuery(form.SortingRequest),
		// TODO add your keyword fields here (cf. userDom example)
	}

	repoView, err := svc.repository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(repoView.Total)
	view.Exemples = make([]CommonView, 0)

	for _, exemple := range repoView.Exemples {
		view.Exemples = append(view.Exemples, svc.fromModelToView(exemple))
	}

	return
}

/* Private methods below */

func (svc *service) fromModelToView(exemple model.Exemple) (view CommonView) {
	view = CommonView{
		Id: exemple.ID,
		// TODO add your custom fields mapping here
	}
	return
}

func (svc *service) fromFormToModel(form CommonForm) (exemple model.Exemple) {
	exemple = model.Exemple{
		// TODO add your custom fields mapping here
	}
	return
}
