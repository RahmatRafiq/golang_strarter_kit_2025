package helpers

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	// if value is empty, return default value
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	// if value is empty, return default value
	if len(value) == 0 {
		return defaultValue
	}

	// convert string to int
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func GetEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)

	// if value is empty, return default value
	if len(value) == 0 {
		return defaultValue
	}

	// convert string to int64
	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return int64Value
}