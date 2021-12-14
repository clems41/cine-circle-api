package main

import (
	"cine-circle-api/cmd/migration-manager/migration"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/sql/sqlConnection"
)

// main will run all migration jobs. These jobs are used to add, update, or delete data from database.
// They can also be used to update your database schema (add or delete column, indexes, etc...).
func main() {
	// Try to connect to PostgresSQL database (using default config from env variables)
	DB, err := sqlConnection.Open(nil)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	defer func() {
		err = sqlConnection.Close(DB)
		if err != nil {
			logger.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		} else {
			logger.Infof("Database connection has been gracefully closed")
		}
	}()

	// Run the main process
	currentVersion, upToDateVersion, err := migration.RunMigrations(DB)

	if err != nil {
		logger.Fatalf("Migrations are NOT OK : %s", err)
	}

	logger.Infof("Data migration process is finished! Jobs are done %d/%d",
		currentVersion, upToDateVersion)
}
