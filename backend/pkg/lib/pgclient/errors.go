package pgclient

import (
	"errors"
)

var ErrConfigParsingFailed = errors.New("pg client: config parsing failed")

var ErrOpeningFailed = errors.New("pg client: opening failed")
