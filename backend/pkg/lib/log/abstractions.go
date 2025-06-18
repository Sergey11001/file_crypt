package log

import (
	"context"
	"log/slog"
)

type Logger interface {
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
}

type LoggerFunc func(ctx context.Context, level slog.Level, msg string, args ...any)

func (f LoggerFunc) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	f(ctx, level, msg, args)
}

type Wither interface {
	With(args ...any) *slog.Logger
}

type WitherFunc func(args ...any) *slog.Logger

func (f WitherFunc) With(args ...any) *slog.Logger {
	return f(args...)
}
