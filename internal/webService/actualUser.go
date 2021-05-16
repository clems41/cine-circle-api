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
	"strconv"
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

// WhoAmI : return username from token
func (auh *actualUserHandler) WhoAmI(req *restful.Request) (user ActualUser, err error) {
	// Get userID from token
	token, err := GetTokenFromAuthenticationHeader(req)
	if err != nil {
		return
	}
	claims, err := CheckToken(token)
	if err != nil {
		return
	}
	subClaims := fmt.Sprintf("%v", claims["sub"])
	userID, err := strconv.Atoi(subClaims)
	if err != nil {
		return user, typedErrors.NewAuthenticationErrorf("cannot find userID from token, got claims = %s", subClaims)
	}

	// Get user's info from DB based on userID (token)
	var userFromDB repositoryModel.User
	err = auh.DB.
		Take(&userFromDB, "id = ?", userID).
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
