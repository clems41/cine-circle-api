package circleDom

import "cine-circle/internal/repository/repositoryModel"

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (view View, err error)
	Update(update Update) (view View, err error)
	Get(get Get) (view View, err error)
	Delete(deletion Deletion) (err error)
	AddUser(updateUser UpdateUser) (view View, err error)
	DeleteUser(updateUser UpdateUser) (view View, err error)
	List(filters Filters) (list ListView, err error)
}

type service struct {
	r repository
}

func NewService(r repository) Service {
	return &service{
		r: r,
	}
}

func (svc *service) Create(creation Creation) (view View, err error) {
	// valid fields
	err = creation.Valid()
	if err != nil {
		return
	}
	// get user that sending request
	userFromRequest, err := svc.r.GetUser(creation.UserIDFromRequest)
	if err != nil {
		return
	}
	// creating circle and adding user in it
	circle := repositoryModel.Circle{
		Name:        creation.Name,
		Description: creation.Description,
		Users: []repositoryModel.User{userFromRequest},
	}
	// save into db
	err = svc.r.Create(&circle)
	if err != nil {
		return
	}
	// convert into view
	view = svc.toView(circle)
	return
}

func (svc *service) Update(update Update) (view View, err error) {
	err = update.Valid()
	if err != nil {
		return
	}
	userFromRequest, err := svc.r.GetUser(update.UserIDFromRequest)
	if err != nil {
		return
	}

	circle, err := svc.r.Get(update.CircleID)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userFromRequest.ID) {
		return view, errNotAuthorized
	}

	if update.Name != "" {
		circle.Name = update.Name
	}
	if update.Description != "" {
		circle.Description = update.Description
	}
	circle.ID = update.CircleID

	err = svc.r.Update(&circle)
	if err != nil {
		return
	}
	view = svc.toView(circle)
	return
}

func (svc *service) Get(get Get) (view View, err error) {
	// get user that sending request
	userFromRequest, err := svc.r.GetUser(get.UserIDFromRequest)
	if err != nil {
		return
	}

	circle, err := svc.r.Get(get.CircleID)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userFromRequest.ID) {
		return view, errNotAuthorized
	}
	view = svc.toView(circle)
	return
}

func (svc *service) Delete(deletion Deletion) (err error) {
	// get user that sending request
	userFromRequest, err := svc.r.GetUser(deletion.UserIDFromRequest)
	if err != nil {
		return
	}

	circle, err := svc.r.Get(deletion.CircleID)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userFromRequest.ID) {
		return errNotAuthorized
	}
	return svc.r.Delete(circle.GetID())
}

func (svc *service) AddUser(updateUser UpdateUser) (view View, err error) {
	// get user that sending request
	userFromRequest, err := svc.r.GetUser(updateUser.UserIDFromRequest)
	if err != nil {
		return
	}

	circle, err := svc.r.Get(updateUser.CircleID)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userFromRequest.ID) {
		return view, errNotAuthorized
	}

	userToAdd, err := svc.r.GetUser(updateUser.UserIDToUpdate)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userToAdd.ID) {
		err = svc.r.AddUserToCircle(userToAdd, &circle)
		if err != nil {
			return
		}
	}

	view = svc.toView(circle)
	return
}

func (svc *service) DeleteUser(updateUser UpdateUser) (view View, err error) {
	// get user that sending request
	userFromRequest, err := svc.r.GetUser(updateUser.UserIDFromRequest)
	if err != nil {
		return
	}

	circle, err := svc.r.Get(updateUser.CircleID)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userFromRequest.ID) {
		return view, errNotAuthorized
	}

	userToDelete, err := svc.r.GetUser(updateUser.UserIDToUpdate)
	if err != nil {
		return
	}

	if !svc.userIsInCircle(circle, userToDelete.ID) {
		return view, errUserNotFound
	}

	err = svc.r.DeleteUserFromCircle(userToDelete.ID, &circle)
	if err != nil {
		return
	}
	view = svc.toView(circle)
	return
}

func (svc *service) List(filters Filters) (list ListView, err error) {
	circles, total, err := svc.r.List(filters)
	if err != nil {
		return
	}
	list.Page = filters.BuildResult(total)
	for _, circle := range circles {
		list.Circles = append(list.Circles, svc.toView(circle))
	}
	return
}

func (svc *service) toView(circle repositoryModel.Circle) (view View) {
	view = View{
		CircleID:    circle.ID,
		Name:        circle.Name,
		Description: circle.Description,
	}
	for _, user := range circle.Users {
		view.Users = append(view.Users, UserView{
			UserID:      user.GetID(),
			Username:    *user.Username,
			DisplayName: user.DisplayName,
		})
	}
	return
}

func (svc *service) userIsInCircle(circle repositoryModel.Circle, userID uint) (isIn bool) {
	for _, userInCircle := range circle.Users {
		if userInCircle.ID == userID {
			return true
		}
	}
	return false
}
