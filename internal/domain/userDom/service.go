package userDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (result Result, err error)
	Update(update Update) (result Result, err error)
	UpdatePassword(updatePassword UpdatePassword) (result Result, err error)
	Delete(delete Delete) (err error)
	Get(get Get) (result Result, err error)
}

type service struct {
	r Repository
}

type Repository interface {
	Create(creation Creation) (result Result, err error)
	Update(update Update) (result Result, err error)
	UpdatePassword(updatePassword UpdatePassword) (result Result, err error)
	Delete(delete Delete) (rr error)
	Get(get Get) (result Result, err error)
	GetHashedPassword(username string) (hashedPassword string, err error)
}

func NewService(r Repository) Service {
	return &service{
		r:                              r,
	}
}

func (svc *service) Create(creation Creation) (result Result, err error) {
	err = creation.Valid()
	if err != nil {
		return
	}
	return svc.r.Create(creation)
}

func (svc *service) Update(update Update) (result Result, err error) {
	err = update.Valid()
	if err != nil {
		return
	}
	return svc.r.Update(update)
}

func (svc *service) UpdatePassword(updatePassword UpdatePassword) (result Result, err error) {
	err = updatePassword.Valid()
	if err != nil {
		return
	}

	user, err := svc.r.Get(Get{UserID: updatePassword.UserID})
	if err != nil {
		return
	}

	hashedPassword, err := svc.r.GetHashedPassword(user.Username)
	if err != nil {
		return
	}

	err = utils.CompareHashAndPassword(hashedPassword, updatePassword.OldPassword)
	if err != nil {
		return result, typedErrors.NewApiBadRequestErrorf(err.Error())
	}

	updatePassword.NewHashedPassword, err = utils.HashPassword(updatePassword.NewPassword, constant.CostHashFunction)
	return svc.r.UpdatePassword(updatePassword)
}

func (svc *service) Delete(delete Delete) (err error) {
	err = delete.Valid()
	if err != nil {
		return
	}
	return svc.r.Delete(delete)
}

func (svc *service) Get(get Get) (result Result, err error) {
	return svc.r.Get(get)
}
