package espresso

import (
	"context"

	"golang.org/x/exp/slog"
)

var (
	defaultLogger = slog.Default()
)

func Debug(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).ErrorContext(ctx, msg, args...)
}

func grabLogger(ctx context.Context) *slog.Logger {
	if rCtx, ok := ctx.(Context); ok {
		return rCtx.Logger()
	}

	return defaultLogger
}
