package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSaltPassword hash and salt password using bcrypt
func HashAndSaltPassword(password string, cost int) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	hashedPassword = string(bytes)
	return
}

// CompareHashAndPassword check if hashedPassword is corresponding to password
func CompareHashAndPassword(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

