package test

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/repository"
	"fmt"
	"gorm.io/gorm"
	"testing"
)

const (
	TemplateDatabaseName = "test_template_db"
	TestingDatabaseNamePrefix = "testing_database_"
)

var DatabaseIndexName = 0

func OpenDatabase(t *testing.T) (DB *gorm.DB, closeFunction func()) {
	// Open connection with real database for creating testing database
	realDB, err := repository.OpenConnection()
	if err != nil {
		t.Fatalf(err.Error())
	}

	testingDatabaseName := fmt.Sprintf("%s%d", TestingDatabaseNamePrefix, DatabaseIndexName)
	DatabaseIndexName++

	err = realDB.
		Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", testingDatabaseName, TemplateDatabaseName)).
		Error
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Using new clean database for testing
	testingDB, err := repository.OpenConnection(testingDatabaseName)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create closeFunction : to be used in testcases
	closeFunction = func() {
		// closing testing database before deleting it from real database
		err = repository.CloseConnection(testingDB)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// delete testing database
		err = realDB.
			Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testingDatabaseName)).
			Error
		if err != nil {
			logger.Sugar.Fatalf(err.Error())
		}

		// closing connection from real database
		err = repository.CloseConnection(realDB)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	return testingDB, closeFunction
}
