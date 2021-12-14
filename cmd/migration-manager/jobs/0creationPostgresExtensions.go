package jobs

import (
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/sliceUtils"
	"fmt"
	"gorm.io/gorm"
)

var extensions = []string{
	"unaccent",
	"fuzzystrmatch",
	"pg_trgm",
	"hstore",
}

type CreationPostgresExtensions struct{}

func (job CreationPostgresExtensions) JobName() string {
	return "#XXXXX - PostgreSQL - Création des extensions Postgres utiles au projet"
}

func (job CreationPostgresExtensions) Job(tx *gorm.DB) (err error) {
	// Create extensions, useful for advanced query
	for _, extension := range extensions {
		err = tx.
			Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", extension)).
			Error
		if err != nil {
			return
		}
	}
	return
}

func (job CreationPostgresExtensions) IsJobDone(DB *gorm.DB) bool {
	list := make([]struct {
		Extname string
	}, 0)
	err := DB.
		Raw("SELECT extname FROM pg_extension").
		Scan(&list).
		Error
	if err != nil {
		logger.Errorf("Error while getting Postgres extensions : %s", err.Error())
		return false
	}

	// Check that all extensions are included in extension list get from database query
	var actualExtensions []string
	for _, elem := range list {
		actualExtensions = append(actualExtensions, elem.Extname)
	}
	if !sliceUtils.StrSlicesContainsAll(extensions, actualExtensions) {
		logger.Errorf("Some extensions are missing. Actual=%s \t Expected=%s", actualExtensions, extensions)
		return false
	}
	return true
}

func (job CreationPostgresExtensions) PreventJobRerun() bool {
	return true
}
