package espresso

import (
	"github.com/googollee/module/log"
)

var (
	LogModule = log.Module
	LogText   = log.TextLogger
	LogJSON   = log.JSONLogger

	DEBUG = log.DEBUG
	WARN  = log.WARN
	INFO  = log.INFO
	ERROR = log.ERROR
)

func logHandling(ctx Context) error {
	method := ctx.Request().Method
	path := ctx.Request().URL.String()

	slog := LogModule.Value(ctx)
	if slog != nil {
		ctx = ctx.WithParent(log.With(ctx, "method", method, "path", path))
	}

	INFO(ctx, "receive http")
	defer INFO(ctx, "finish http")

	ctx.Next()

	return nil
}
