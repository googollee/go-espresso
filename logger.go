package espresso

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

var (
	defaultLogger = slog.Default()
)

func WithLog(logger *slog.Logger) ServerOption {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

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

func MiddlewareLogRequest(ctx Context) error {
	req := ctx.Request()
	ctx = WithLogAttr(ctx, "span", time.Now().Unix(), "method", req.Method, "path", req.URL.Path)

	respW := &responseWriter{
		ResponseWriter: ctx.ResponseWriter(),
		logger:         ctx.Logger(),
	}
	ctx = WithResponseWriter(ctx, respW)

	Info(ctx, "Request")
	defer func() {
		if err := ctx.Error(); err != nil {
			Error(ctx, "Error", "error", err.Error())
		}
		if p := recover(); p != nil {
			Error(ctx, "Panic", "panic", fmt.Sprintf("%s", p))
		}
		Info(ctx, "Response", "code", respW.responseCode)
	}()

	ctx.Next()

	return nil
}
