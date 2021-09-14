package test

import (
	"cine-circle/pkg/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestMainContainsAllDomains(t *testing.T) {
	// read the whole file main.go to check if all domain have been imported inside
	bytes, err := ioutil.ReadFile("../../cmd/cine-circle-api/main.go")
	if err != nil {
		t.Fatalf(err.Error())
	}
	mainFile := string(bytes)

	// List all domains in order to check them in main file
	files, err := ioutil.ReadDir("../../internal/domain")
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			// Check that package is imported
			require.True(t, strings.Contains(mainFile, file.Name()), "main.go should contains domain %s", file.Name())
			// Check that webservice has been created
			creationStr := file.Name() + ".NewHandler("
			require.True(t, strings.Contains(mainFile, creationStr), "main.go should create handler for %s", file.Name())
		}
	}
}

func TestRepositoriesContainsAllDomains(t *testing.T) {
	exceptions := []string{"rootDom"}
	// read the whole file main.go to check if all domain have been imported inside
	bytes, err := ioutil.ReadFile("../../internal/repository/repositories.go")
	if err != nil {
		t.Fatalf(err.Error())
	}
	mainFile := string(bytes)

	// List all domains in order to check them in main file
	files, err := ioutil.ReadDir("../../internal/domain")
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			// Skip check for exceptions
			if utils.SliceContainsStr(exceptions, file.Name()) {
				continue
			}
			// Check that package is imported
			require.True(t, strings.Contains(mainFile, file.Name()), "repositories.go should contains domain %s", file.Name())
			// Check that repository has been created
			creationStr := file.Name() + ".NewRepository("
			require.True(t, strings.Contains(mainFile, creationStr), "repositories.go should create repository for %s", file.Name())
		}
	}
}
