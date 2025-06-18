package pgclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func New(config Config, logger Logger) (*pgxpool.Pool, error) {
	if logger == nil {
		panic("pg client: nil logger")
	}

	u, err := NewURL(config)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigParsingFailed, err)
	}
	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(
			func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
				args := make([]any, 0, len(data))
				for k, v := range data {
					args = append(args, slog.Any(k, v))
				}
			},
		),
		LogLevel: tracelog.LogLevelTrace,
	}

	c, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpeningFailed, err)
	}

	return c, nil
}
