package log

import (
	"context"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func (l *Logger) CheckHealth(ctx context.Context) error {
	return nil
}

type Option func(l *Logger) error

func OptionSlogger(logger *slog.Logger) Option {
	if logger == nil {
		logger = slog.Default()
	}

	return func(l *Logger) error {
		l.logger = logger
		return nil
	}
}
