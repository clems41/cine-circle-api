package migration

import (
	"cine-circle-api/internal/repository/postgresRepositories"
	"cine-circle-api/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Manager struct {
	gorm.Model
	Version        int
	LastRunJobName string
}

func (Manager) TableName() string {
	return "migration_manager"
}

func RunMigrations(DB *gorm.DB) (currentVersion int, upToDateVersion int, err error) {

	migrations := GetMigrationJobs()

	upToDateVersion = len(migrations)

	/* Retrieving current state of the DB */
	err = DB.AutoMigrate(&Manager{})
	if err != nil {
		return
	}
	var currentState Manager
	err = DB.FirstOrCreate(&currentState).Error
	if err != nil {
		return currentVersion, upToDateVersion, errors.WithStack(err)
	}

	// Migration jobs start here
	// They cannot access internal project code : that makes them re-runnable in any database state with
	// any codebase version.
	// That requires to redefine types etc. for each job in order for them to be independent
	// Also : all these jobs execute within a single transaction, which should avoid any heart attack
	// while deploying to production
	logger.Infof("We begin the new job system, so we start a transaction")
	tx := DB.Begin()
	rollback := func(err error) {
		logger.Errorf("Something went wrong: rolling back : %+v", err)
		tx.Rollback()
	}

	// First migrate all new tables before inserting any new data
	err = postgresRepositories.Migrate(tx)
	if err != nil {
		rollback(err)
		return currentState.Version, upToDateVersion, err
	}

	err = applyMigrationList(tx, migrations, &currentState)
	if err != nil {
		rollback(err)
		return currentState.Version, upToDateVersion, err
	}

	if currentState.Version != upToDateVersion {
		err = errors.New(fmt.Sprintf("Migrations are not up to date, current version = %d but should be %d", currentState.Version, upToDateVersion))
		rollback(err)
		return currentState.Version, upToDateVersion, err
	}

	logger.Infof("Everything went OK. Commiting the changes")
	tx.Commit()
	return currentState.Version, upToDateVersion, nil
}

func applyMigrationList(tx *gorm.DB, migrationList []Migration, currentState *Manager) (err error) {
	/* Iterate over jobs, run next relevant job and increment DB version */
	for jobNumber, migration := range migrationList {
		if jobNumber == currentState.Version {

			logger.Infof("Manager is running '%s'", migration.JobName())

			if migration.PreventJobRerun() && migration.IsJobDone(tx) {

				err = errors.New(
					fmt.Sprintf("Manager failed to run '%s' because the job seems to be already done. Please check the state of the database.", migration.JobName()),
				)

				return
			}

			err = migration.Job(tx)

			if err == nil {

				if migration.IsJobDone(tx) {

					currentState.Version += 1
					currentState.LastRunJobName = migration.JobName()

					err = tx.Save(&currentState).Error
					if err != nil {
						return
					}

					logger.Infof("Manager has successfully run '%s'", migration.JobName())

				} else {

					currentState.LastRunJobName = migration.JobName()

					tx.Save(&currentState)

					err = errors.New(
						fmt.Sprintf("Manager failed to run '%s'. There were no error, but the job seems to be incomplete. Please check the state of the database.", migration.JobName()),
					)

					return
				}

			} else { // If Job sent an error

				logger.Errorf("Manager failed to run '%s' : %+v", migration.JobName(), err)

				return
			}

		} else {

			logger.Infof("Skipping already run job %d : '%s'", jobNumber, migration.JobName())
		}
	}

	return
}
