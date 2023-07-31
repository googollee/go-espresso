package espresso

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

var (
	loggerKey     = InjectKey("espresso.log.Log")
	defaultLogger = slog.Default()
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

type logger struct {
	slogger *slog.Logger
}

type LoggerOption func(*logger)

func Logger(option ...LoggerOption) HandleFunc {
	logger := logger{
		slogger: slog.Default(),
	}

	for _, op := range option {
		op(&logger)
	}

	return logger.handle
}

func LogWithSlog(log *slog.Logger) LoggerOption {
	return func(l *logger) {
		l.slogger = log
	}
}

func (l *logger) handle(ctx Context) error {
	req := ctx.Request()

	args := []any{"span", time.Now().Unix()}
	args = append(args, "method", req.Method)
	args = append(args, "path", req.URL.Path)

	logger := l.slogger.With(args)
	ctx.InjectValue(loggerKey, logger)

	Info(ctx, "received")

	defer func() {
		if err := ctx.Err(); err != nil {
			code := http.StatusInternalServerError
			if hc, ok := err.(HTTPCoder); ok {
				code = hc.HTTPCode()
			}
			Error(ctx, "done with error", "code", code, "error", err)
			return
		}
		Info(ctx, "done")
	}()

	ctx.Next()

	return nil
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
