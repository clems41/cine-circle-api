package watchlistDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
	Save(elem repositoryModel.Watchlist) (err error)
	Delete(elem repositoryModel.Watchlist) (err error)
	Find(elem *repositoryModel.Watchlist) (err error)
	List(filters Filters) (list []repositoryModel.Watchlist, total int64, err error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {

	err := r.DB.AutoMigrate(&repositoryModel.Watchlist{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating watchlistRepository : %s", err.Error())
	}

}

func (r *Repository) Save(elem repositoryModel.Watchlist) (err error) {
	var movie repositoryModel.Movie
	err = r.DB.
		Take(&movie, "id = ?", elem.MovieID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.DB.Create(&elem).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (r *Repository) Delete(elem repositoryModel.Watchlist) (err error) {
	var movie repositoryModel.Movie
	err = r.DB.
		Take(&movie, "id = ?", elem.MovieID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.DB.
		Delete(&repositoryModel.Watchlist{}, "user_id = ? AND movie_id = ?", elem.UserID, elem.MovieID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (r *Repository) Find(elem *repositoryModel.Watchlist) (err error) {
	var movie repositoryModel.Movie
	err = r.DB.
		Take(&movie, "id = ?", elem.MovieID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.DB.Take(elem).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (r *Repository) List(filters Filters) (list []repositoryModel.Watchlist, total int64, err error) {
	query := r.DB.
		Preload("Movie").
		Preload("User").
		Order("created_at desc").
		Where("user_id = ?", filters.UserID)

	if filters.PageSize != 0 {
		query = query.Limit(filters.PageSize)
	}

	query = query.Offset(filters.Offset())

	err = query.
		Find(&list).
		Limit(-1).
		Offset(-1).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return
}
