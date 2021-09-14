package utils

import "os"

// GetDefaultOrFromEnv : return Value from envVariable name, if not defines return defaultValue
func GetDefaultOrFromEnv(defaultValue, fromEnv string) string {
	envVal := os.Getenv(fromEnv)
	if envVal != "" {
		return envVal
	} else {
		return defaultValue
	}
}
