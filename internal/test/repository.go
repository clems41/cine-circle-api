package test

import (
	"cine-circle/internal/repository"
	"cine-circle/internal/utils"
	"gorm.io/gorm"
	"testing"
)

const (
	EnvTestDatabaseName = "TEST_DATABASE_NAME"
	DefaultTestDatabaseName = "test_cine_circle"
)

var (
	testDatabaseName = utils.GetDefaultOrFromEnv(DefaultTestDatabaseName, EnvTestDatabaseName)
)

func OpenDatabase(t *testing.T) (DB *gorm.DB, closeFunction func()) {
	// Open connection
	DB, err := repository.OpenConnection(testDatabaseName)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// create closeFunction : to be used in testcases
	closeFunction = func() {
		err = repository.CloseConnection(DB)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	return
}
