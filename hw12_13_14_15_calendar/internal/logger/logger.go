package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string) (*Logger, error) {
	var slogLevel slog.Level
	if err := slogLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("cannot parse logger level: %w", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	}))
	slog.SetDefault(logger)
	return &Logger{
		logger: logger,
	}, nil
}

func (l Logger) Debug(_ context.Context, msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l Logger) Info(_ context.Context, msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l Logger) Warn(_ context.Context, msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l Logger) Error(_ context.Context, err error, msg string, args ...any) {
	l.logger.Error(msg, append(args, "error", err)...)
}
