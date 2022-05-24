package pgRepositories

import (
	"gorm.io/gorm"
)

type PgRepository interface {
	Migrate() (err error)
}

// Migrate will call all migration functions defined in repository/model.
// To add new migration function, please add it into list function below.
func Migrate(tx *gorm.DB) (err error) {
	for _, repository := range repositories(tx) {
		err = repository.Migrate()
		if err != nil {
			return
		}
	}
	return
}

// list return all migration functions that should be called during database migration.
// Please add here your new migration method.
func repositories(tx *gorm.DB) (repositories []PgRepository) {
	repositories = append(repositories, NewCircle(tx))
	repositories = append(repositories, NewMedia(tx))
	repositories = append(repositories, NewRecommendation(tx))
	repositories = append(repositories, NewUser(tx))
	return
}
