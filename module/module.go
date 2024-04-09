// Package module provides a way to do dependency injection, with type-safe, without performance penalty.
// See examples for the basic usage.
package module

import (
	"context"
	"reflect"
)

// BuildFunc is the constructor of an Instance.
type BuildFunc[T any] func(context.Context) (T, error)

type moduleKey string

// Provider is the interface to provide an Instance.
type Provider interface {
	key() moduleKey
	value(ctx context.Context) (any, error)
}

// Module provides a module to inject an instance with its type.
// As Module implements Provider interface, a Module instance could be added into a Repo instance.
type Module[T any] struct {
	Provider
	moduleKey moduleKey
	builder   BuildFunc[T]
}

// New creates a new module with type `T` and the constructor `builder`.
func New[T any](builder BuildFunc[T]) Module[T] {
	var t T
	return Module[T]{
		moduleKey: moduleKey(reflect.TypeOf(&t).Elem().String()),
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

	return v.(T)
}

func (m Module[T]) key() moduleKey {
	return m.moduleKey
}

func (m Module[T]) value(ctx context.Context) (any, error) {
	return m.builder(ctx)
}
