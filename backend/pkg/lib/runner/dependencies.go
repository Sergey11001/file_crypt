package runner

// Logger is implemented by the [slog.Logger] type.
type Logger interface {
	Debug(msg string, args ...any)
}
