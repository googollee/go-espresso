package espresso

import (
	"context"
	"log/slog"
	"os"

	"github.com/googollee/go-espresso/module"
)

var (
	LogModule = module.New[*slog.Logger]()
	LogText   = LogModule.ProvideWithFunc(func(ctx context.Context) (*slog.Logger, error) {
		handler := slog.NewTextHandler(os.Stderr, nil)
		return slog.New(handler), nil
	})
	LogJSON = LogModule.ProvideWithFunc(func(ctx context.Context) (*slog.Logger, error) {
		handler := slog.NewJSONHandler(os.Stderr, nil)
		return slog.New(handler), nil
	})
)

func DEBUG(ctx context.Context, msg string, args ...any) {
	logger := LogModule.Value(ctx)
	if logger == nil {
		return
	}

	logger.DebugContext(ctx, msg, args...)
}

func INFO(ctx context.Context, msg string, args ...any) {
	logger := LogModule.Value(ctx)
	if logger == nil {
		return
	}

	logger.InfoContext(ctx, msg, args...)
}

func WARN(ctx context.Context, msg string, args ...any) {
	logger := LogModule.Value(ctx)
	if logger == nil {
		return
	}

	logger.WarnContext(ctx, msg, args...)
}

func ERROR(ctx context.Context, msg string, args ...any) {
	logger := LogModule.Value(ctx)
	if logger == nil {
		return
	}

	logger.ErrorContext(ctx, msg, args...)
}

func logHandling(ctx Context) error {
	method := ctx.Request().Method
	path := ctx.Request().URL.String()

	INFO(ctx, "receive http", "method", method, "path", path)
	defer INFO(ctx, "finish http", "method", method, "path", path)

	ctx.Next()

	return nil
}
