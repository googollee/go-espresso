package module

import (
	"context"
)

type errBuildPanic struct {
	error
}

type BuildFunc[T Instance] func(ctx context.Context) (T, error)

type Builder interface {
	Key
	build(ctx context.Context) (Instance, error)
}

type builder[T Instance] struct {
	key
	buildFn BuildFunc[T]
}

func (b builder[T]) build(ctx context.Context) (Instance, error) {
	return b.buildFn(ctx)
}
