package service

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
	"cine-circle/internal/repository"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	cost = 8
	expirationDuration = 1 * 24 * time.Hour
	secretTokenEnv = "SECRET_TOKEN"
	secretTokenDefault = "secret"
	tokenKind = "Bearer"
	tokenHeader = "Authorization"
)

func HashAndSaltPassword(password string, user *model.User) typedErrors.CustomError {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return typedErrors.NewApiUnprocessableEntityError(err)
	}
	user.Hash = string(hashedPassword)
	user.Password = ""
	return typedErrors.NoErr
}

func passwordIsCorrect(hashPassword, password string) (typedErrors.CustomError, bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return typedErrors.NewApiBadCredentialsError(err), false
	}
	return typedErrors.NoErr, true
}

func getHashAndPassword(auth string) (typedErrors.CustomError, string, string, string) {
	result := strings.Split(auth, " ")
	if len(result) != 2 {
		return typedErrors.ErrApiUnprocessableEntity, "", "", ""
	}
	loginPasswordEncoded := result[1]
	loginPasswordDecoded, err := base64.StdEncoding.DecodeString(loginPasswordEncoded)
	if err != nil {
		return typedErrors.ErrApiUnprocessableEntity, "", "", ""
	}
	pair := strings.Split(string(loginPasswordDecoded), ":")
	if len(result) != 2 {
		return typedErrors.ErrApiUnprocessableEntity, "", "", ""
	}
	username := pair[0]
	password := pair[1]
	db, err2 := repository.OpenConnection()
	if err2.IsNotNil() {
		return err2, "", "", ""
	}
	defer db.Close()
	var hashedPassword string
	res := db.DB().Table("users").Select("hash").Find(&hashedPassword, "username = ?", username)
	if res.Error != nil {
		return typedErrors.NewRepositoryQueryFailedError(res.Error), "", "", ""
	}
	if res.RowsAffected != 1 || hashedPassword == "" {
		return typedErrors.ErrApiUserBadCredentials, "", "", ""
	}
	return typedErrors.NoErr, hashedPassword, password, username
}

func GetTokenFromAuthentication(auth string) (typedErrors.CustomError, string) {
	err, hashPassword, password, username := getHashAndPassword(auth)
	if err.IsNotNil() {
		return err, ""
	}
	err2, same := passwordIsCorrect(hashPassword, password)
	if err2.IsNotNil() {
		return err2, ""
	}
	if !same {
		return typedErrors.ErrApiUserBadCredentials, ""
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,    jwt.MapClaims{
		"iss": "cine-circle",
		"sub":  username,
		"aud": "any",
		"exp": time.Now().Add(expirationDuration).Unix(),
	})

	jwtToken, _:= token.SignedString([]byte(utils.GetDefaultOrFromEnv(secretTokenDefault, secretTokenEnv)))
	return typedErrors.NoErr, jwtToken
}

func CheckTokenAndGetUsername(req *restful.Request) (typedErrors.CustomError, string) {
	token := req.HeaderParameter(tokenHeader)
	if token == "" {
		return typedErrors.ErrApiUserBadCredentials, ""
	}
	res := strings.Split(token, " ")
	if len(res) != 2 {
		return typedErrors.ErrApiUserBadCredentials, ""
	}
	if res[0] != tokenKind {
		return typedErrors.ErrApiUserBadCredentials, ""
	}
	tokenStr := res[1]
	claims := jwt.MapClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetDefaultOrFromEnv(secretTokenDefault, secretTokenEnv)), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return typedErrors.NewApiBadCredentialsError(err), ""
		}
	}
	if !tkn.Valid {
		logger.Sugar.Debugf("Error while getting token : Token not valid")
		return typedErrors.ErrApiUserBadCredentials, ""
	}
	return typedErrors.NoErr, fmt.Sprintf("%v", claims["sub"])
}
