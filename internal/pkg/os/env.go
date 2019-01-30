package os

import "os"

func GetEnvOrDefault(key, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return value
}
