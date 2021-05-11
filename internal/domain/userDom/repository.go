package userDom

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/pkg/logger"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
	Create(user *repositoryModel.User) (err error)
	Save(user *repositoryModel.User) (err error)
	Delete(userID domain.IDType) (err error)
	Get(get Get) (user repositoryModel.User, err error)
	Search(filters Filters) (result []repositoryModel.User, err error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {

	err := r.DB.AutoMigrate(&repositoryModel.User{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating repository : %s", err.Error())
	}

	err = r.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username_display_name ON users USING GIST ((username || display_name) gist_trgm_ops)").Error
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when creating index : %s", err.Error())
	}

}

func (r *Repository) Create(user *repositoryModel.User) (err error) {
	return r.DB.
		Create(&user).
		Error
}

func (r *Repository) Save(user *repositoryModel.User) (err error) {
	return r.DB.
		Save(&user).
		Error
}

func (r *Repository) Delete(userID domain.IDType) (rr error) {
	return r.DB.
		Delete(&repositoryModel.User{}, "id = ?", userID).
		Error
}

func (r *Repository) Get(get Get) (user repositoryModel.User, err error) {
	return r.getUser(get)
}

func (r *Repository) Search(filters Filters) (users []repositoryModel.User, err error) {

	keyword := "%" + filters.Keyword + "%"

	err = r.DB.
		Where("concat(username || display_name) ILIKE ?", keyword).
		Find(&users).
		Error

	return
}

func (r *Repository) getUser(get Get) (user repositoryModel.User, err error) {
	query := r.DB

	if get.UserID != 0 {
		query = query.Where("id = ?", get.UserID)
	}
	if get.Username != "" {
		query = query.Where("username = ?", get.Username)
	}
	if get.Email != "" {
		query = query.Where("email = ?", get.Email)
	}

	err = query.
		Take(&user).
		Error

	return
}
