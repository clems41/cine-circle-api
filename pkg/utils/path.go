package utils

import (
	"cine-circle-api/internal/constant/settingsConst"
	"cine-circle-api/pkg/utils/pathUtils"
)

// GetResourcePath return path of resources directory, should work for all environments (CI, local, master, etc...)
func GetResourcePath() (resourcesPath string, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}
	resourcesPath = rootPath + settingsConst.RelativeResourcePathFromRootDir

	return
}
