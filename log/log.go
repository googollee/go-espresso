package log

import (
	"context"
	"time"

	"github.com/googollee/go-espresso"
	"golang.org/x/exp/slog"
)

func Debug(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).DebugCtx(ctx, msg, args)
}

func Info(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).InfoCtx(ctx, msg, args)
}

func Warn(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).WarnCtx(ctx, msg, args)
}

func Error(ctx context.Context, msg string, args ...any) {
	grabLogger(ctx).ErrorCtx(ctx, msg, args)
}

func With(logger *slog.Logger) espresso.HandleFunc {
	defaultLogger = logger
	return func(ctx espresso.Context) error {
		req := ctx.Request()

		id := time.Now().Unix()
		// espresso.Context provides espresso.Injecting.
		// WithArgs should inject into ctx without creating a new instance.
		WithArgs(ctx, "method", req.Method, "path", req.URL.Path, "id", id)

		Info(ctx, "start")
		defer func() {
			Info(ctx, "end")
		}()

		ctx.Next()

		return nil
	}
}

func WithArgs(ctx context.Context, args ...any) context.Context {
	logger := grabLogger(ctx).With(args...)

	if injector, ok := ctx.(espresso.Injecting); ok {
		injector.InjectValue(loggerKey, logger)
		return ctx
	}

	return context.WithValue(ctx, loggerKeyType(loggerKey), logger)
}

type loggerKeyType string

var (
	loggerKey     = espresso.InjectKey("espresso.log.Log")
	defaultLogger = slog.Default()
)

func grabLogger(ctx context.Context) *slog.Logger {
	v := ctx.Value(loggerKey)
	if v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}

	return defaultLogger
}
