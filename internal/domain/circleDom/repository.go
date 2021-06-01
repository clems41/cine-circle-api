package circleDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
	Create(circle *repositoryModel.Circle) (err error)
	Update(circle *repositoryModel.Circle) (err error)
	Get(circleID uint) (circle repositoryModel.Circle, err error)
	Delete(circleID uint) (err error)
	GetUser(userID uint) (user repositoryModel.User, err error)
	AddUserToCircle(user repositoryModel.User, circle *repositoryModel.Circle) (err error)
	DeleteUserFromCircle(userID uint, circle *repositoryModel.Circle) (err error)
	List(filters Filters) (circles []repositoryModel.Circle, total int64, err error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {

	err := r.DB.AutoMigrate(&repositoryModel.Circle{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating circleRepository : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_circle_user_circle ON circle_user (circle_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_circle_user_user ON circle_user (user_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

}

func (r *Repository) Create(circle *repositoryModel.Circle) (err error) {
	err = r.DB.
		Create(circle).
		Error
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) Update(circle *repositoryModel.Circle) (err error) {
	err = r.DB.
		Save(circle).
		Error
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) Get(circleID uint) (circle repositoryModel.Circle, err error) {
	err = r.DB.
		Preload("Users").
		Take(&circle, "id = ?", circleID).
		Error
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) GetUser(userID uint) (user repositoryModel.User, err error) {
	err = r.DB.
		Take(&user, "id = ?", userID).
		Error
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) Delete(circleID uint) (err error) {
	err = r.DB.
		Delete(&repositoryModel.Circle{}, "id = ?", circleID).
		Error
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) AddUserToCircle(user repositoryModel.User, circle *repositoryModel.Circle) (err error) {
	err = r.DB.
		Model(circle).
		Association("Users").
		Append([]repositoryModel.User{user})
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) DeleteUserFromCircle(userID uint, circle *repositoryModel.Circle) (err error) {
	err = r.DB.
		Model(circle).
		Association("Users").
		Delete([]repositoryModel.User{
			{
				Metadata: repositoryModel.Metadata{ID: userID},
			},
		})
	if err != nil {
		errors.WithStack(err)
	}
	return
}

func (r *Repository) List(filters Filters) (circles []repositoryModel.Circle, total int64, err error) {
	query := r.DB.
		Preload("Users").
		Where("id IN (SELECT circle_id FROM circle_user WHERE user_id = ?)", filters.UserID)

	if filters.PageSize != 0 {
		query = query.Limit(filters.PageSize)
	}

	query = query.Offset(filters.Offset())

	err = query.
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
