package repository

import (
	"cine-circle-api/internal/repository/model"
	"gorm.io/gorm"
)

// MigrationFunction define who migration methods should be implemented
type MigrationFunction func(tx *gorm.DB) error

// Migrate will call all migration functions defined in repository/model.
// To add new migration function, please add it into list function below.
func Migrate(tx *gorm.DB) (err error) {
	for _, migrationFunc := range list() {
		err = migrationFunc(tx)
		if err != nil {
			return
		}
	}
	return
}

// list return all migration functions that should be called during database migration.
// Please add here your new migration method.
func list() (migrationList []MigrationFunction) {
	migrationList = append(migrationList, model.MigrateUser)
	migrationList = append(migrationList, model.MigrateMedia)
	migrationList = append(migrationList, model.MigrateCircle)
	migrationList = append(migrationList, model.MigrateRecommendation)
	// TODO add your new repository migration method here
	return
}
