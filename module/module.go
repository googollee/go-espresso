package module

import (
	"context"
	"fmt"
)

// Instance holds a instance which is injected with the module.
type Instance interface {
	// CheckHealth returns `nil` if the instance is health.
	CheckHealth(ctx context.Context) error
}

// BuildFunc is the constructor of an Instance.
type BuildFunc[T Instance] func(context.Context) (T, error)

type moduleKey string

// Provider is the interface to provide an Instance.
type Provider interface {
	key() moduleKey
	value(ctx context.Context) (Instance, error)
}

// Module provides a module to inject an instance with its type.
// As Module implements Provider interface, a Module instance could be added into a Repo instance.
type Module[T Instance] struct {
	Provider
	moduleKey moduleKey
	builder   BuildFunc[T]
}

// New creates a new module with type `T` and the constructor `builder`.
func New[T Instance](builder BuildFunc[T]) Module[T] {
	var t T
	return Module[T]{
		moduleKey: moduleKey(fmt.Sprintf("%T", t)),
		builder:   builder,
	}
}

// Value returns an instance of T which is injected to the context.
func (m Module[T]) Value(ctx context.Context) T {
	var null T

	v := ctx.Value(m.moduleKey)
	if v == nil {
		return null
	}

	ret, ok := v.(T)
	if !ok {
		return null
	}

	return ret
}

func (m Module[T]) key() moduleKey {
	return m.moduleKey
}

func (m Module[T]) value(ctx context.Context) (Instance, error) {
	return m.builder(ctx)
}
