package log

import (
	"context"
	"log/slog"

	"github.com/googollee/go-espresso/basetype"
)

func getLogger(ctx context.Context) *slog.Logger {
	v := ctx.Value(ctxKey)
	if v == nil {
		return slog.Default()
	}

	ret, ok := v.(*Logger)
	if !ok {
		return slog.Default()
	}

	return ret.logger
}

func WithAttr(ctx basetype.Context, args ...any) basetype.Context {
	v := ctx.Value(ctxKey)
	if v == nil {
		return ctx
	}

	orgCtx, ok := v.(*Logger)
	if !ok {
		return ctx
	}

	retCtx := &Logger{
		logger: orgCtx.logger.With(args...),
	}

	return ctx.WithParent(context.WithValue(ctx, ctxKey, retCtx))
}

func Debug(ctx context.Context, msg string, args ...any) {
	logger := getLogger(ctx)
	logger.DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	logger := getLogger(ctx)
	logger.InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	logger := getLogger(ctx)
	logger.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	logger := getLogger(ctx)
	logger.ErrorContext(ctx, msg, args...)
}
