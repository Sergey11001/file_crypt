package log

import (
	"fmt"
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Level string `default:"info"`
}

func New(config LoggerConfig) (*slog.Logger, error) {
	var level slog.Level
	switch config.Level {
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidLevel, config.Level)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})), nil
}
