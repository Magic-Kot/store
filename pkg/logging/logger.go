package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	LogLevel string
}

func NewLogger(ctx context.Context) (context.Context, error) {
	//zerolog.SetGlobalLevel(zerolog.TraceLevel)

	logger := zerolog.New(os.Stdout)

	ctx = logger.WithContext(ctx)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return ctx, nil
}
