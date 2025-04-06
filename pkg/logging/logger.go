package logging

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type LoggerDeps struct {
	Level string
}

func NewLogger(cfg *LoggerDeps) (*zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}
	zerolog.SetGlobalLevel(logLevel)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger, nil
}
