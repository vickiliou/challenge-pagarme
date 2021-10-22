package config

import (
	"os"
)

const (
	APIKeyEnv      = "PAGARMEAPI_APIKEY"
	APIEndpointEnv = "PAGARMEAPI_ENDPOINT"
)

// GetEnvVar returns the value of the environment variable and error
func GetEnvVar(envVar string) (string, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return "", errMissingEnvVar
	}
	return value, nil
}
