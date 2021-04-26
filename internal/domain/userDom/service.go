package userDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (result Result, err error)
	Update(update Update) (result Result, err error)
	UpdatePassword(updatePassword UpdatePassword) (result Result, err error)
	Delete(delete Delete) (err error)
	Get(get Get) (result Result, err error)
	Search(filters Filters) (result []Result, err error)
	GenerateTokenFromAuthenticationHeader(header string) (token string, err error)
	getUsernameAndPasswordFromAuthenticationHeader(header string) (username string, password string, err error)
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
	Search(filters Filters) (result []Result, err error)
	GetHashedPassword(get Get) (hashedPassword string, err error)
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
	// Hash and salt password
	hashedPassword, err := utils.HashAndSaltPassword(creation.Password, constant.CostHashFunction)
	if err != nil {
		return
	}
	// Save hashed and salt password as user's password
	creation.Password = hashedPassword
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

	hashedPassword, err := svc.r.GetHashedPassword(Get{UserID: updatePassword.UserID})
	if err != nil {
		return
	}

	err = utils.CompareHashAndPassword(hashedPassword, updatePassword.OldPassword)
	if err != nil {
		return result, typedErrors.NewApiBadRequestError(err)
	}

	updatePassword.NewHashedPassword, err = utils.HashAndSaltPassword(updatePassword.NewPassword, constant.CostHashFunction)
	if err != nil {
		return result, typedErrors.NewApiBadRequestError(err)
	}
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
	err = get.Valid()
	if err != nil {
		return
	}
	return svc.r.Get(get)
}

func (svc *service) Search(filters Filters) (result []Result, err error) {
	err = filters.Valid()
	if err != nil {
		return
	}
	return svc.r.Search(filters)
}

func (svc *service) GenerateTokenFromAuthenticationHeader(header string) (token string, err error) {
	username, password, err := svc.getUsernameAndPasswordFromAuthenticationHeader(header)
	if err != nil {
		return
	}

	hashedPassword, err := svc.r.GetHashedPassword(Get{Username: username})
	if err != nil {
		return
	}

	err = utils.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		err = typedErrors.NewApiBadCredentialsError(err)
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": constant.IssToken,
		"sub": username,
		"aud": "any",
		"exp": time.Now().Add(constant.ExpirationDuration).Unix(),
	})

	return jwtToken.SignedString([]byte(constant.TokenKey))
}

func (svc *service) getUsernameAndPasswordFromAuthenticationHeader(header string) (username string, password string, err error) {
	result := strings.Split(header, " ")
	if len(result) != 2 {
		err = typedErrors.NewApiBadCredentialsErrorf("Header format is not correct")
		return
	}
	loginPasswordEncoded := result[1]
	loginPasswordDecoded, err := base64.StdEncoding.DecodeString(loginPasswordEncoded)
	if err != nil {
		err = typedErrors.NewApiBadCredentialsError(err)
		return
	}
	pair := strings.Split(string(loginPasswordDecoded), ":")
	if len(result) != 2 {
		err = typedErrors.NewApiBadCredentialsErrorf("Encoded login:password is not correct")
		return
	}
	username = pair[0]
	password = pair[1]
	return
}
