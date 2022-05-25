package postgresRepositories

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository.Recommendation = (*recommendationPgRepository)(nil)
var _ PgRepository = (*recommendationPgRepository)(nil)

type recommendationPgRepository struct {
	DB *gorm.DB
}

func NewRecommendation(DB *gorm.DB) *recommendationPgRepository {
	return &recommendationPgRepository{DB: DB}
}

func (repo *recommendationPgRepository) Migrate() (err error) {
	err = repo.DB.
		AutoMigrate(&model.Recommendation{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (repo *recommendationPgRepository) Create(recommendation *model.Recommendation) (err error) {
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

func (repo *recommendationPgRepository) Search(form repository.RecommendationSearchForm) (recommendations []model.Recommendation, total int64, err error) {
	query := repo.DB.
		Preload("Movie").
		Preload("Circles").
		Preload("Circles.Users").
		Preload("Sender")

	if form.MediaId != 0 {
		query = query.Where("media_id = ?", form.MediaId)
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
		Find(&recommendations).
		Limit(-1).
		Offset(-1).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return
}
