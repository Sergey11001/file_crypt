package config

import (
	"errors"
)

var ErrEnvLoadingFailed = errors.New("config: env loading failed")

var ErrEnvParsingFailed = errors.New("config: env parsing failed")
