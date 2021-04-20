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