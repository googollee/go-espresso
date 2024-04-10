package module

import "context"

// Provider is the interface to provide an Instance.
type Provider interface {
	key() moduleKey
	value(ctx context.Context) (any, error)
}

type funcProvider[T any] struct {
	moduleKey moduleKey
	ctor      BuildFunc[T]
}

// ProvideWithFunc returns a provider which provides instances creating from `ctor` function.
func (m *Module[T]) ProvideWithFunc(ctor BuildFunc[T]) Provider {
	return &funcProvider[T]{
		moduleKey: m.moduleKey,
		ctor:      ctor,
	}
}

func (p funcProvider[T]) key() moduleKey {
	return p.moduleKey
}

func (p funcProvider[T]) value(ctx context.Context) (any, error) {
	return p.ctor(ctx)
}
