package movieDom

import (
	"cine-circle/internal/repository/repositoryModel"
	logger "cine-circle/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
	Save(movie *repositoryModel.Movie) (err error)
	Get(movieId uint) (movie repositoryModel.Movie, err error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {

	err := r.DB.AutoMigrate(&repositoryModel.Movie{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating movieRepository : %s", err.Error())
	}

}

func (r *Repository) Save(movie *repositoryModel.Movie) (err error) {
	err = r.DB.
		Save(movie).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (r *Repository) Get(movieId uint) (movie repositoryModel.Movie, err error) {
	err = r.DB.
		Take(&movie, "id = ?", movieId).
		Error
	if err != nil {
		return movie, errors.WithStack(err)
	}
	return
}
