package utils

import (
	"cine-circle-api/internal/constant/settingsConst"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// TestGetResourcePath check that GetResourcePath is working by opening test.txt file from cmd/migration-manager/resources directory
func TestGetResourcePath(t *testing.T) {
	err := os.Setenv("APPLICATION_NAME", "cine-circle-api")
	require.NoError(t, err)
	// Test without INIT_ROOT_PROJECT_PATH env variable
	filename := "test.txt"
	resourcesPath, err := GetResourcePath()
	require.NoError(t, err)
	_, err = os.Open(resourcesPath + filename)
	require.NoError(t, err)

	// Test with INIT_ROOT_PROJECT_PATH env variable
	actualDir, err := os.Getwd()
	require.NoError(t, err)
	initRootProjectPath := actualDir + "/../.."
	err = os.Setenv(settingsConst.InitRootProjectPathEnv, initRootProjectPath)
	require.NoError(t, err)
	resourcesPath, err = GetResourcePath()
	require.NoError(t, err)
	_, err = os.Open(resourcesPath + filename)
	require.NoError(t, err)
}
