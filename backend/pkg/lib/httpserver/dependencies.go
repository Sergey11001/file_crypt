package httpserver

import (
	"univer/pkg/lib/log"
)

// Logger is implemented by the [slog.Logger] type.
type Logger interface {
	log.Logger

	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}
