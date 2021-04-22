package main

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/repository"
	"cine-circle/internal/test"
	"cine-circle/internal/utils"
	"fmt"
)

var (
	testDatabaseName = utils.GetDefaultOrFromEnv(test.DefaultTestDatabaseName, test.EnvTestDatabaseName)
)

func main() {
	logger.InitLogger()
	// Try to connect to PostgresSQL repository
	DB, err := repository.OpenConnection()
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}
	closeFunction := func() {
		err = repository.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}

	// Remove all connections from old test database
	err = DB.
		Exec(fmt.Sprintf("select pg_terminate_backend(pid) from pg_stat_activity where datname='%s';", testDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Delete old testing database
	err = DB.
		Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Create new testing database
	err = DB.
		Exec(fmt.Sprintf("CREATE DATABASE %s", testDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Close and open connection on new testing database
	closeFunction()
	DB, err = repository.OpenConnection(testDatabaseName)
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}
	defer func() {
		err = repository.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}()

	// Migrate all tables into testing database
	repositories := repository.NewAllRepositories(DB)
	repositories.Migrate()

/*	// Make sure the template is not modified.
	err = DB.Exec(fmt.Sprintf("ALTER DATABASE %s WITH ALLOW_CONNECTIONS 0;", testDatabaseName)).Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}*/
}
