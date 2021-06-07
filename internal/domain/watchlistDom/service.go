package watchlistDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Service = (*service)(nil)

type Service interface {
	AddMovie(creation Creation) (err error)
	DeleteMovie(deletion Delete) (err error)
	AlreadyExists(check Check) (exists bool, err error)
	List(filters Filters) (list List, err error)
}

type service struct {
	r repository
}

func NewService(r repository) Service {
	return &service{
		r: r,
	}
}

func (service *service) AddMovie(creation Creation) (err error) {
	elem := repositoryModel.Watchlist{
		UserID:   creation.UserID,
		MovieID:  creation.MovieID,
	}
	err = service.r.Save(elem)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errMovieNotFound
		}
	}
	return
}

func (service *service) DeleteMovie(deletion Delete) (err error) {
	elem := repositoryModel.Watchlist{
		UserID:   deletion.UserID,
		MovieID:  deletion.MovieID,
	}
	err = service.r.Delete(elem)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errMovieNotFound
		}
	}
	return
}

func (service *service) AlreadyExists(check Check) (exists bool, err error) {
	elem := repositoryModel.Watchlist{
		UserID:   check.UserID,
		MovieID:  check.MovieID,
	}
	err = service.r.Find(&elem)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (service *service) List(filters Filters) (list List, err error) {
	watchlist, total, err := service.r.List(filters)
	if err != nil {
		return
	}
	list.Page = filters.PaginationRequest.BuildResult(total)
	for _, element := range watchlist {
		if element.Movie != nil {
			list.Movies = append(list.Movies, MovieView{
				ID:          element.Movie.GetID(),
				Title:       element.Movie.Title,
				PosterPath:  element.Movie.PosterPath,
				ReleaseDate: element.Movie.ReleaseDate,
			})
		}
	}
	return
}
