package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func New[T any]() (T, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		var zero T

		return zero, fmt.Errorf("%w: %w", ErrEnvLoadingFailed, err)
	}

	var config T

	err = envconfig.Process("", &config)
	if err != nil {
		var zero T

		return zero, fmt.Errorf("%w: %w", ErrEnvParsingFailed, err)
	}

	return config, nil
}
