package repository

import (
	"cine-circle/internal/domain/authenticationDom"
	"gorm.io/gorm"
)

var _ authenticationDom.Repository = (*authenticationRepository)(nil)

type authenticationRepository struct {
	DB *gorm.DB
}

func NewAuthenticationRepository(DB *gorm.DB) *authenticationRepository {
	return &authenticationRepository{DB: DB}
}

func (r authenticationRepository) Migrate() {

}

func (r authenticationRepository) GetHashedPassword(username string) (hashedPassword string, err error) {
	var user User
	err = r.DB.
		Select("hashed_password").
		Find(&user, "username = ?", username).
		Error
	hashedPassword = user.HashedPassword
	return
}