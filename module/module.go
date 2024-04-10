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

// Module provides a module to inject and retreive an instance with its type.
type Module[T any] struct {
	moduleKey moduleKey
}

// New creates a new module with type `T` and the constructor `builder`.
func New[T any]() Module[T] {
	var t T
	return Module[T]{
		moduleKey: moduleKey(reflect.TypeOf(&t).Elem().String()),
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
