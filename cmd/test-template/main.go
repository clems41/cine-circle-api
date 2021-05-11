package main

import (
	"cine-circle/internal/repository"
	"cine-circle/internal/repository/postgres"
	"cine-circle/internal/test"
	logger "cine-circle/pkg/logger"
	"fmt"
)

func main() {
	logger.InitLogger()
	// Try to connect to PostgresSQL repository
	DB, err := postgres.OpenConnection()
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}
	closeFunction := func() {
		err = postgres.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}

	// Remove all connections from old test database
	err = DB.
		Exec(fmt.Sprintf("select pg_terminate_backend(pid) from pg_stat_activity where datname='%s';", test.TemplateDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Delete old testing database
	err = DB.
		Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", test.TemplateDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Create new testing database
	err = DB.
		Exec(fmt.Sprintf("CREATE DATABASE %s", test.TemplateDatabaseName)).
		Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}

	// Close and open connection on new testing database
	closeFunction()
	DB, err = postgres.OpenConnection(test.TemplateDatabaseName)
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}
	defer func() {
		err = postgres.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}()

	// Migrate all tables into testing database
	repositories := repository.NewAllRepositories(DB)
	repositories.Migrate()

	// Make sure the template is not modified.
/*	err = DB.Exec(fmt.Sprintf("ALTER DATABASE %s WITH ALLOW_CONNECTIONS 0;", test.TemplateDatabaseName)).Error
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}*/
}
