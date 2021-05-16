package userDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"encoding/base64"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (view View, err error)
	Update(update  Update) (view View, err error)
	UpdatePassword(updatePassword UpdatePassword) (view View, err error)
	Delete(delete Delete) (err error)
	Get(get Get) (view View, err error)
	GetOwnInfo(userID uint) (view ViewMe, err error)
	Search(filters Filters) (views []View, err error)
	GenerateTokenFromAuthenticationHeader(header string) (token string, err error)
	getUsernameAndPasswordFromAuthenticationHeader(header string) (username string, password string, err error)
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
	// Validate fields
	err = creation.Valid()
	if err != nil {
		return
	}
	// Hash and salt password
	hashedPassword, err := utils.HashAndSaltPassword(creation.Password, constant.CostHashFunction)
	if err != nil {
		return
	}
	// Save hashed and salt password as user's password
	creation.Password = hashedPassword

	username := strings.ToLower(creation.Username)

	user := repositoryModel.User{
		Username:       &username,
		DisplayName:    creation.DisplayName,
		Email:          creation.Email,
		HashedPassword: creation.Password,
	}

	err = svc.r.Create(&user)
	if err != nil {
		return
	}

	view = svc.toView(user)
	return
}

func (svc *service) Update(update Update) (view View, err error) {
	// Validate fields
	err = update.Valid()
	if err != nil {
		return
	}

	// Get old user from DB
	user, err := svc.r.Get(Get{UserID: update.UserID})
	if err != nil {
		return
	}

	// Update specific fields
	user.DisplayName = update.DisplayName
	user.Email = update.Email

	// Save new user info
	err = svc.r.Save(&user)
	if err != nil {
		return
	}

	view = svc.toView(user)
	return
}

func (svc *service) UpdatePassword(updatePassword UpdatePassword) (view View, err error) {
	// Validate fields
	err = updatePassword.Valid()
	if err != nil {
		return
	}

	// Get old user from DB
	user, err := svc.r.Get(Get{UserID: updatePassword.UserID})
	if err != nil {
		return
	}

	err = utils.CompareHashAndPassword(user.HashedPassword, updatePassword.OldPassword)
	if err != nil {
		return view, errBadLoginPassword
	}

	newHashedPassword, err := utils.HashAndSaltPassword(updatePassword.NewPassword, constant.CostHashFunction)
	if err != nil {
		return view, err
	}
	// Save new user info
	user.HashedPassword = newHashedPassword
	err = svc.r.Save(&user)
	if err != nil {
		return
	}

	view = svc.toView(user)
	return
}

func (svc *service) Delete(delete Delete) (err error) {
	// Validate fields
	err = delete.Valid()
	if err != nil {
		return
	}

	return svc.r.Delete(delete.UserID)
}

func (svc *service) Get(get Get) (view View, err error) {
	// Validate fields
	err = get.Valid()
	if err != nil {
		return
	}
	user, err := svc.r.Get(get)
	if err != nil {
		return
	}

	view = svc.toView(user)
	return
}

func (svc *service) GetOwnInfo(userID uint) (view ViewMe, err error) {
	user, err := svc.r.Get(Get{UserID: userID})
	if err != nil {
		return
	}

	view = svc.toViewMe(user)
	return
}

func (svc *service) Search(filters Filters) (views []View, err error) {
	err = filters.Valid()
	if err != nil {
		return
	}
	users, err := svc.r.Search(filters)
	if err != nil {
		return
	}

	for _, user := range users {
		views = append(views, svc.toView(user))
	}
	return
}

func (svc *service) GenerateTokenFromAuthenticationHeader(header string) (token string, err error) {
	username, password, err := svc.getUsernameAndPasswordFromAuthenticationHeader(header)
	if err != nil {
		return
	}

	user, err := svc.r.Get(Get{Username: username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return token, errBadLoginPassword
		}
		return
	}

	err = utils.CompareHashAndPassword(user.HashedPassword, password)
	if err != nil {
		return token, errBadLoginPassword
	}

	return utils.GenerateTokenWithUserID(user.GetID())
}

func (svc *service) getUsernameAndPasswordFromAuthenticationHeader(header string) (username string, password string, err error) {
	result := strings.Split(header, constant.BearerTokenDelimiterForHeader)
	if len(result) != 2 {
		err = typedErrors.NewAuthenticationErrorf("Header format is not correct")
		return
	}
	loginPasswordEncoded := result[1]
	loginPasswordDecoded, err := base64.StdEncoding.DecodeString(loginPasswordEncoded)
	if err != nil {
		err = typedErrors.NewAuthenticationErrorf(err.Error())
		return
	}
	pair := strings.Split(string(loginPasswordDecoded), constant.UsernamePasswordDelimiterForHeader)
	if len(result) != 2 {
		err = typedErrors.NewAuthenticationErrorf("Encoded login:password is not correct")
		return
	}
	username = pair[0]
	password = pair[1]
	return
}

func (svc *service) toView(user repositoryModel.User) (view View) {
	view = View {
		UserID:      user.GetID(),
		DisplayName: user.DisplayName,
	}
	if user.Username != nil {
		view.Username = *user.Username
	}
	return
}

func (svc *service) toViewMe(user repositoryModel.User) (view ViewMe) {
	view = ViewMe {
		UserID:      user.GetID(),
		DisplayName: user.DisplayName,
	}
	if user.Username != nil {
		view.Username = *user.Username
	}
	view.Email = user.Email
	return
}
