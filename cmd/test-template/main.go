package main

import (
	"cine-circle-api/internal/repository/postgres/pgRepositories"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/sql/sqlTest"
)

// main will create template database and migrate all repositories to create database schemas, tables, columns, indexes, etc...
// This database template is used for testing. Create new database and migrate all repositories at the beginning of each test will take too much time.
// So, we are creating this template before running all testcases, and for each test, new testing database will be created from this template (without doing any migrations).
func main() {
	err := sqlTest.CreateTestTemplateDatabase(pgRepositories.Migrate, logger.Logger())
	if err != nil {
		logger.Fatalf(err.Error())
	} else {
		logger.Infof("Template database has been correctly created")
	}
}
