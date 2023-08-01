package espresso

import (
	"context"

	"golang.org/x/exp/slog"
)

var (
	loggerKey     = InjectKey("espresso.log.Log")
	defaultLogger = slog.Default()
)

func Debug(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).DebugCtx(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).InfoCtx(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).WarnCtx(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).ErrorCtx(ctx, msg, args...)
}

func grabLogger(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}

	return defaultLogger
}
