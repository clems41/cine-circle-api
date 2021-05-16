package circleDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/pkg/logger"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
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
