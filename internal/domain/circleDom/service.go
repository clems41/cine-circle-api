package circleDom

import (
	"cine-circle-api/internal/repository/instance/circleRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/sql/gormUtils"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(form CreateForm) (view CreateView, err error)
	AddUser(form AddUserForm) (view AddUserView, err error)
	DeleteUser(form DeleteUserForm) (view DeleteUserView, err error)
	Update(form UpdateForm) (view UpdateView, err error)
	Get(form GetForm) (view GetView, err error)
	Delete(form DeleteForm) (err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	repository circleRepository.Repository
}

func NewService(repository circleRepository.Repository) Service {
	return &service{
		repository: repository,
	}
}

func (svc *service) Create(form CreateForm) (view CreateView, err error) {
	circle := svc.fromFormToModel(form.CommonForm)
	err = svc.repository.Create(&circle)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Update(form UpdateForm) (view UpdateView, err error) {
	circle, ok, err := svc.repository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}
	circle = svc.fromFormToModel(form.CommonForm)
	circle.ID = form.CircleId
	err = svc.repository.Save(&circle)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	circle, ok, err := svc.repository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}
	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Delete(form DeleteForm) (err error) {
	_, ok, err := svc.repository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return errCircleNotFound
	}

	err = svc.repository.Delete(form.CircleId)
	if err != nil {
		return
	}
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	repoForm := circleRepository.SearchForm{
		PaginationQuery: gormUtils.PaginationQuery{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		CircleName: form.CircleName,
	}

	repoView, err := svc.repository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(repoView.Total)
	view.Circles = make([]CommonView, 0)

	for _, circle := range repoView.Circles {
		view.Circles = append(view.Circles, svc.fromModelToView(circle))
	}

	return
}

func (svc *service) AddUser(form AddUserForm) (view AddUserView, err error) {
	return
}

func (svc *service) DeleteUser(form DeleteUserForm) (view DeleteUserView, err error) {
	return
}

/* Private methods below */

func (svc *service) fromModelToView(circle model.Circle) (view CommonView) {
	view = CommonView{
		Id:          circle.ID,
		Name:        circle.Name,
		Description: circle.Description,
	}
	for _, user := range circle.Users {
		view.Users = append(view.Users, UserView{
			Id:        user.ID,
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Username:  user.Username,
		})
	}
	return
}

func (svc *service) fromFormToModel(form CommonForm) (circle model.Circle) {
	circle = model.Circle{
		Name:        form.Name,
		Description: form.Description,
	}
	return
}
