package module

import (
	"context"
	"fmt"
)

type Instance interface {
	CheckHealth(ctx context.Context) error
}

type Key interface {
	name() key
}

type key string

func (k key) name() key {
	return k
}

func (k key) String() string {
	return fmt.Sprintf("module(%s)", string(k))
}

type Module[T Instance] struct {
	key
}

func New[T Instance]() Module[T] {
	var t T
	return Module[T]{
		key: key(fmt.Sprintf("%T", t)),
	}
}

func (m Module[T]) Value(ctx context.Context) T {
	var zero T

	v := ctx.Value(m.key)
	if v == nil {
		return zero
	}

	ret, ok := v.(T)
	if !ok {
		return zero
	}

	return ret
}

func (m Module[T]) Builder(fn BuildFunc[T]) Builder {
	return &builder[T]{
		key:     m.key,
		buildFn: fn,
	}
}
