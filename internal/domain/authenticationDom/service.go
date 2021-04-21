package authenticationDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	GenerateTokenFromAuthenticationHeader(header string) (token string, err error)
	Create(creation userDom.Creation) (result userDom.Result, err error)
	getUsernameAndPasswordFromAuthenticationHeader(header string) (username string, password string, err error)
}

type service struct {
	r Repository
	userRepository userDom.Repository
}

type Repository interface {
}

func NewService(r Repository, userRepository userDom.Repository) Service {
	return &service{
		r:                              r,
		userRepository: userRepository,
	}
}

func (svc *service) GenerateTokenFromAuthenticationHeader(header string) (token string, err error) {
	username, password, err := svc.getUsernameAndPasswordFromAuthenticationHeader(header)
	if err != nil {
		return
	}

	hashedPassword, err := svc.userRepository.GetHashedPassword(userDom.Get{Username: username})
	if err != nil {
		return
	}

	err = utils.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		err = typedErrors.NewApiBadCredentialsErrorf(err.Error())
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,    jwt.MapClaims{
		"iss": constant.IssToken,
		"sub":  username,
		"aud": "any",
		"exp": time.Now().Add(constant.ExpirationDuration).Unix(),
	})

	return jwtToken.SignedString([]byte(constant.TokenKey))
}

func (svc *service) Create(creation userDom.Creation) (result userDom.Result, err error) {
	// Hash and salt password
	hashedPassword, err := utils.HashPassword(creation.Password, constant.CostHashFunction)
	if err != nil {
		return
	}
	// Save hashed and salt password as user's password
	creation.Password = string(hashedPassword)

	return svc.userRepository.Create(creation)
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
		err = typedErrors.NewApiBadCredentialsErrorf(err.Error())
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