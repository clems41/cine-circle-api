package repository

import (
	"cine-circle/internal/domain/userDom"
	"gorm.io/gorm"
)

var _ userDom.Repository = (*userRepository)(nil)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *userRepository {
	return &userRepository{DB: DB}
}

func (r userRepository) Migrate() {

}