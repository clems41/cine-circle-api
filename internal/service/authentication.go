package service

import (
	"cine-circle/internal/database"
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
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

func HashAndSaltPassword(password string, user *model.User) model.CustomError {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return model.NewCustomError(err, model.ErrInternalApiUnprocessableEntity.HttpCode(), model.ErrInternalApiUnprocessableEntityCode)
	}
	user.Hash = string(hashedPassword)
	user.Password = ""
	return model.NoErr
}

func passwordIsCorrect(hashPassword, password string) (model.CustomError, bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return model.NewCustomError(err, model.ErrInternalApiUserBadCredentials.HttpCode(), model.ErrInternalApiUserBadCredentialsCode), false
	}
	return model.NoErr, true
}

func getHashAndPassword(auth string) (model.CustomError, string, string, string) {
	result := strings.Split(auth, " ")
	if len(result) != 2 {
		return model.ErrInternalApiUnprocessableEntity, "", "", ""
	}
	loginPasswordEncoded := result[1]
	loginPasswordDecoded, err := base64.StdEncoding.DecodeString(loginPasswordEncoded)
	if err != nil {
		return model.ErrInternalApiUnprocessableEntity, "", "", ""
	}
	pair := strings.Split(string(loginPasswordDecoded), ":")
	if len(result) != 2 {
		return model.ErrInternalApiUnprocessableEntity, "", "", ""
	}
	username := pair[0]
	password := pair[1]
	db, err2 := database.OpenConnection()
	if err2.IsNotNil() {
		return err2, "", "", ""
	}
	defer db.Close()
	var hashedPassword string
	res := db.DB().Table("users").Select("hash").Find(&hashedPassword, "username = ?", username)
	if res.Error != nil {
		return model.NewCustomError(res.Error, model.ErrInternalDatabaseQueryFailed.HttpCode(), model.ErrInternalDatabaseQueryFailedCode), "", "", ""
	}
	if res.RowsAffected != 1 || hashedPassword == "" {
		return model.ErrInternalApiUserBadCredentials, "", "", ""
	}
	return model.NoErr, hashedPassword, password, username
}

func GetTokenFromAuthentication(auth string) (model.CustomError, string) {
	err, hashPassword, password, username := getHashAndPassword(auth)
	if err.IsNotNil() {
		return err, ""
	}
	err2, same := passwordIsCorrect(hashPassword, password)
	if err2.IsNotNil() {
		return err2, ""
	}
	if !same {
		return model.ErrInternalApiUserBadCredentials, ""
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,    jwt.MapClaims{
		"iss": "cine-circle",
		"sub":  username,
		"aud": "any",
		"exp": time.Now().Add(expirationDuration).Unix(),
	})

	jwtToken, _:= token.SignedString([]byte(utils.GetDefaultOrFromEnv(secretTokenDefault, secretTokenEnv)))
	return model.NoErr, jwtToken
}

func CheckTokenAndGetUsername(req *restful.Request) (model.CustomError, string) {
	token := req.HeaderParameter(tokenHeader)
	if token == "" {
		return model.ErrInternalApiUserBadCredentials, ""
	}
	res := strings.Split(token, " ")
	if len(res) != 2 {
		return model.ErrInternalApiUserBadCredentials, ""
	}
	if res[0] != tokenKind {
		return model.ErrInternalApiUserBadCredentials, ""
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
			return model.NewCustomError(err, model.ErrInternalApiUserBadCredentials.HttpCode(), model.ErrInternalApiUserBadCredentialsCode), ""
		}
	}
	if !tkn.Valid {
		logger.Sugar.Debugf("Error while getting token : Token not valid")
		return model.ErrInternalApiUserBadCredentials, ""
	}
	return model.NoErr, fmt.Sprintf("%v", claims["sub"])
}
