package circleDom

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"cine-circle-api/pkg/utils/searchUtils"
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
	circleRepository repository.Circle
	userRepository   repository.User
}

func NewService(circleRepository repository.Circle, userRepository repository.User) Service {
	return &service{
		circleRepository: circleRepository,
		userRepository:   userRepository,
	}
}

func (svc *service) Create(form CreateForm) (view CreateView, err error) {
	circle := model.Circle{
		Name:        form.Name,
		Description: form.Description,
	}
	err = svc.circleRepository.Save(&circle)
	if err != nil {
		return
	}

	err = svc.circleRepository.AddUser(form.UserId, &circle)
	if err != nil {
		return
	}

	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Update(form UpdateForm) (view UpdateView, err error) {
	circle, ok, err := svc.circleRepository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}

	found := svc.userIsInCircle(form.UserId, circle)
	if !found {
		return view, errCircleNotFound
	}

	circle.Name = form.Name
	circle.Description = form.Description
	circle.ID = form.CircleId
	err = svc.circleRepository.Save(&circle)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	circle, ok, err := svc.circleRepository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}

	found := svc.userIsInCircle(form.UserId, circle)
	if !found {
		return view, errCircleNotFound
	}

	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) Delete(form DeleteForm) (err error) {
	circle, ok, err := svc.circleRepository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return errCircleNotFound
	}

	found := svc.userIsInCircle(form.UserId, circle)
	if !found {
		return errCircleNotFound
	}

	err = svc.circleRepository.Delete(form.CircleId)
	if err != nil {
		return
	}
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	repoForm := repository.CircleSearchForm{
		PaginationRequest: searchUtils.PaginationRequest{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		UserId: form.UserId,
	}

	circles, total, err := svc.circleRepository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(total)
	view.Circles = make([]CommonView, 0)

	for _, circle := range circles {
		view.Circles = append(view.Circles, svc.fromModelToView(circle))
	}

	return
}

func (svc *service) AddUser(form AddUserForm) (view AddUserView, err error) {
	circle, ok, err := svc.circleRepository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}

	found := svc.userIsInCircle(form.UserId, circle)
	if !found {
		return view, errCircleNotFound
	}

	_, ok, err = svc.userRepository.Get(form.UserIdToAdd)
	if err != nil {
		return
	}
	if !ok {
		return view, errUserNotFound
	}

	err = svc.circleRepository.AddUser(form.UserIdToAdd, &circle)
	view.CommonView = svc.fromModelToView(circle)
	return
}

func (svc *service) DeleteUser(form DeleteUserForm) (view DeleteUserView, err error) {
	circle, ok, err := svc.circleRepository.Get(form.CircleId)
	if err != nil {
		return
	}
	if !ok {
		return view, errCircleNotFound
	}

	found := svc.userIsInCircle(form.UserId, circle)
	if !found {
		return view, errCircleNotFound
	}

	_, ok, err = svc.userRepository.Get(form.UserIdToDelete)
	if err != nil {
		return
	}
	if !ok {
		return view, errUserNotFound
	}

	err = svc.circleRepository.DeleteUser(form.UserIdToDelete, &circle)
	view.CommonView = svc.fromModelToView(circle)
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

func (svc *service) userIsInCircle(userId uint, circle model.Circle) (isIn bool) {
	isIn = false
	for _, circleUser := range circle.Users {
		if userId == circleUser.ID {
			isIn = true
		}
	}
	return
}
