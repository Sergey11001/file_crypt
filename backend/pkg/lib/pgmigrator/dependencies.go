package pgmigrator

// Logger is implemented by the [slog.Logger] type.
type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

// Project is implemented by the [app.Project] type.
type Project interface {
	Name() string
}
