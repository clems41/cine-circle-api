package postgresRepositories

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository.Circle = (*circlePgRepository)(nil)
var _ PgRepository = (*circlePgRepository)(nil)

type circlePgRepository struct {
	DB *gorm.DB
}

func NewCircle(DB *gorm.DB) *circlePgRepository {
	return &circlePgRepository{DB: DB}
}

func (repo *circlePgRepository) Migrate() (err error) {
	err = repo.DB.
		AutoMigrate(&model.Circle{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (repo *circlePgRepository) Get(circleId uint) (circle model.Circle, ok bool, err error) {
	err = repo.DB.
		Preload("Users").
		Take(&circle, circleId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return circle, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *circlePgRepository) Save(circle *model.Circle) (err error) {
	// TODO create if not exists
	err = repo.DB.
		Omit("Users").
		Save(circle).
		Error
	return errors.WithStack(err)
}

func (repo *circlePgRepository) Delete(circleId uint) (err error) {
	err = repo.DB.
		Exec("DELETE FROM circle_users where circle_id = ?", circleId).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Delete(&model.Circle{}, circleId).
		Error
	return errors.WithStack(err)
}

func (repo *circlePgRepository) Search(form repository.CircleSearchForm) (circles []model.Circle, total int64, err error) {
	query := repo.DB

	err = query.
		Where("id IN (SELECT circle_id FROM circle_users WHERE user_id = ?)", form.UserId).
		Offset(form.Offset()).
		Limit(form.PageSize).
		Find(&circles).
		Limit(-1).
		Offset(-1).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return
}

func (repo *circlePgRepository) AddUser(userId uint, circle *model.Circle) (err error) {
	err = repo.DB.
		Exec("INSERT INTO circle_users(user_id, circle_id) VALUES (?, ?)", userId, circle.ID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Preload("Users").
		Take(circle).
		Error
	return errors.WithStack(err)
}

func (repo *circlePgRepository) DeleteUser(userId uint, circle *model.Circle) (err error) {
	err = repo.DB.
		Exec("DELETE FROM circle_users where user_id = ? and circle_id = ?", userId, circle.ID).
		Error
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DB.
		Preload("Users").
		Take(circle).
		Error
	return errors.WithStack(err)
}
