package contextx

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

type contextKeyLogger struct{}

func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger{}, logger)
}

func LoggerFromContextOrDefault(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(contextKeyLogger{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return logger
}
