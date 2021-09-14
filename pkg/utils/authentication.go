package utils

import (
	"cine-circle/internal/constant"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// HashAndSaltPassword hash and salt password using bcrypt
func HashAndSaltPassword(password string, cost int) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return hashedPassword, errors.WithStack(err)
	}
	hashedPassword = string(bytes)
	return
}

// CompareHashAndPassword check if hashedPassword is corresponding to password
func CompareHashAndPassword(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateTokenWithUsername(username string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":               constant.IssToken,
		constant.UserClaims: username,
		"aud":               "any",
		"exp":               time.Now().Add(constant.ExpirationDuration).Unix(),
	})

	tokenKey := GetDefaultOrFromEnv(constant.SecretTokenDefault, constant.SecretTokenEnv)

	return jwtToken.SignedString([]byte(tokenKey))
}
