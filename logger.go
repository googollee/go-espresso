package espresso

import (
	"context"
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
	slogger      *slog.Logger
	withMethod   bool
	withPath     bool
	withArgs     []any
	startMessage func(Context) string
	endMessage   func(Context) string
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

func LogWithArgs(args ...any) LoggerOption {
	return func(l *logger) {
		l.withArgs = args
	}
}

func LogWithMethod() LoggerOption {
	return func(l *logger) {
		l.withMethod = true
	}
}

func LogWithPath() LoggerOption {
	return func(l *logger) {
		l.withPath = true
	}
}

func LogWithMessage(start, end func(Context) string) LoggerOption {
	return func(l *logger) {
		l.startMessage = start
		l.endMessage = end
	}
}

func (l *logger) handle(ctx Context) error {
	req := ctx.Request()

	args := []any{"span", time.Now().Unix()}
	if l.withMethod {
		args = append(args, "method", req.Method)
	}
	if l.withPath {
		args = append(args, "path", req.URL.Path)
	}
	if len(l.withArgs) > 0 {
		args = append(args, l.withArgs...)
	}

	logger := l.slogger.With(args)
	ctx.InjectValue(loggerKey, logger)

	if l.startMessage != nil {
		Info(ctx, l.startMessage(ctx))
	}

	defer func() {
		if l.endMessage != nil {
			Info(ctx, l.endMessage(ctx))
		}
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
