package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, cost int) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	hashedPassword = string(bytes)
	return
}

func CompareHashAndPassword(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

