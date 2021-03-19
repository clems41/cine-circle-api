package utils

import "os"

func GetDefaultOrFromEnv(defaultValue, fromEnv string) string {
	envVal := os.Getenv(fromEnv)
	if envVal != "" {
		return envVal
	} else {
		return defaultValue
	}
}
