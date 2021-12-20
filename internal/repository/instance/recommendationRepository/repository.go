package recommendationRepository

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/internal/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	Create(recommendation *model.Recommendation) (err error)
	Search(repoForm SearchForm) (repoView SearchView, err error)
}

type repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *repository {
	return &repository{DB: DB}
}

func (repo *repository) Create(recommendation *model.Recommendation) (err error) {
	err = repo.DB.
		Create(recommendation).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Preload("Movie").
		Preload("Circles").
		Preload("Circles.Users").
		Preload("Sender").
		Take(recommendation).
		Error
	return errors.WithStack(err)
}

func (repo *repository) Search(form SearchForm) (view SearchView, err error) {
	query := repo.DB.
		Preload("Movie").
		Preload("Circles").
		Preload("Circles.Users").
		Preload("Sender")

	if form.MovieId != 0 {
		query = query.Where("movie_id = ?", form.MovieId)
	}
	switch form.Type {
	case recommendationConst.SentType:
		query = query.Where("sender_id = ?", form.UserId)
	case recommendationConst.ReceivedType:
		query = query.Where("id IN (SELECT recommendation_id FROM recommendation_circles WHERE circle_id IN ("+
			"SELECT circle_id FROM circle_users WHERE user_id = ?))", form.UserId)
	case recommendationConst.AllType:
		query = query.Where("(id IN (SELECT recommendation_id FROM recommendation_circles WHERE circle_id IN ("+
			"SELECT circle_id FROM circle_users WHERE user_id = ?))) OR sender_id = ?", form.UserId, form.UserId)
	}

	err = query.
		Offset(form.Offset()).
		Limit(form.PageSize).
		Find(&view.Recommendations).
		Limit(-1).
		Offset(-1).
		Count(&view.Total).
		Error

	if err != nil {
		return view, errors.WithStack(err)
	}
	return
}
