package setupTestCase

import (
	"cine-circle-api/pkg/sql/sqlTest"
	"gorm.io/gorm"
	"os"
	"testing"
)

// OpenCleanDatabaseFromTemplate will set environment variables APPLICATION_NAME (needed it to avoid adding anv variables for each test manually)
// Then create connection to new testing database using setupTestCase.OpenCleanDatabaseFromTemplate
func OpenCleanDatabaseFromTemplate(t *testing.T) (DB *gorm.DB, closeFunction func()) {
	// Set APPLICATION_NAME env variable
	// Avoid to set env variables for all tests
	err := os.Setenv(envApplicationName, applicationName)
	if err != nil {
		t.Fatalf(err.Error())
	}
	DB, closeFunction = sqlTest.OpenCleanDatabaseFromTemplate(t)
	return
}
