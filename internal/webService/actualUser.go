package webService

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

var (
	ActualUserHandler *actualUserHandler
)

type actualUserHandler struct {
	DB *gorm.DB
}

func NewActualUserHandler(DB *gorm.DB) *actualUserHandler {
	return &actualUserHandler{DB: DB}
}

type ActualUser struct {
	ID       uint
	Username string
	Roles    []string
}

// GetUserFromClaims : return username from token
func (auh *actualUserHandler) GetUserFromClaims(claims jwt.MapClaims) (user ActualUser, err error) {
	username := fmt.Sprintf("%v", claims[constant.UserClaims])

	// Get user's info from DB based on userID (token)
	var userFromDB repositoryModel.User
	err = auh.DB.
		Take(&userFromDB, "username = ?", username).
		Error
	if err != nil {
		return user, errors.WithStack(err)
	}

	// Return only useful user's info
	user = ActualUser{
		ID:       userFromDB.GetID(),
		Username: *userFromDB.Username,
		Roles:    nil,
	}
	// TODO add roles into user's info
	return
}

// GetTokenFromAuthenticationHeader : return token in header request
func GetTokenFromAuthenticationHeader(req *restful.Request) (token string, err error) {
	tokenHeader := req.HeaderParameter(constant.TokenHeader)
	if tokenHeader == "" {
		return token, typedErrors.NewAuthenticationErrorf("cannot find header %s", constant.TokenHeader)
	}
	res := strings.Split(tokenHeader, constant.BearerTokenDelimiterForHeader)
	if len(res) != 2 {
		return token, typedErrors.NewAuthenticationErrorf("cannot find token from header %s", constant.TokenHeader)
	}
	if res[0] != constant.TokenKind {
		return token, typedErrors.NewAuthenticationErrorf("token is not typed as %s", constant.TokenKind)
	}
	token = res[1]
	return
}

// CheckToken : check validity of token from header request and return claims
func CheckToken(token string) (claims jwt.MapClaims, err error) {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tokenKey := utils.GetDefaultOrFromEnv(constant.SecretTokenDefault, constant.SecretTokenEnv)
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, typedErrors.NewAuthenticationErrorf(err.Error())
		}
	}
	if tkn == nil || !tkn.Valid {
		return claims, typedErrors.NewAuthenticationErrorf("token is not valid")
	}
	return
}
