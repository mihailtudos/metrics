package utils

import (
	"os"
	"strconv"
)

func OverrideStringEnvValueWithOsEnv(field *string, envName string) {
	if envValue, ok := os.LookupEnv(envName); ok {
		*field = envValue
	}
}

func OverrideIntEnvValueWithOsEnv(field *int, envName string) error {
	if envValue, ok := os.LookupEnv(envName); ok {
		val, err := strconv.Atoi(envValue)
		if err != nil {
			return err
		}

		*field = val
	}

	return nil
}
