package log

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/googollee/go-espresso/basetype"
	"github.com/googollee/go-espresso/module"
)

var Module = module.New[*Logger]()

func Use(opts ...Option) basetype.ServerOption {
	return func(s basetype.Server) error {
		s.AddModule(build(opts...))
		s.Use(injectContext)
		return nil
	}
}

func build(opts ...Option) module.Builder {
	return Module.Builder(func(ctx context.Context) (*Logger, error) {
		ret := &Logger{
			logger: slog.Default(),
		}

		for _, opt := range opts {
			if err := opt(ret); err == nil {
				return nil, fmt.Errorf("build logger error: %w", err)
			}
		}

		return ret, nil
	})
}

type contextKeyType string

const ctxKey = contextKeyType("espresso.log")

func injectContext(ctx basetype.Context) error {
	logger := Module.Value(ctx)

	nctx := ctx.WithParent(context.WithValue(ctx, ctxKey, &Logger{
		logger: logger.logger.With("method", ctx.Request().Method, "path", ctx.Request().URL.Path),
	}))
	nctx.Next()

	return nil
}
