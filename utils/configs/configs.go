package utils

import (
	"os"
	"strconv"
)

func OverrideStringEnvValueWithOsEnv[T any](field *T, envName string) error {
	if envValue, ok := os.LookupEnv(envName); ok {
		switch v := any(field).(type) {
			case *string:
				*v = envValue
			case *int:
				val, err := strconv.Atoi(envValue)
				if err != nil {
					return err
				}
				*v = val
			case *bool:
				val, err := strconv.ParseBool(envValue)
				if err != nil {
					return err
				}
				
				*v = val
			default:
				panic("Unsupported type")
		}
	}

	return nil
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
