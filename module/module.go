package module

import (
	"context"
	"fmt"
	"reflect"
)

type Instance interface {
	CheckHealth(context.Context) error
}

type BuildFunc[T Instance] func(context.Context) (T, error)

type ModuleKey interface {
	contextKey() contextKey
	build(*buildContext) error
}
type contextKey string

type Module[T Instance] struct {
	ModuleKey
	name    contextKey
	buildFn BuildFunc[T]
}

func New[T Instance](buildFunc BuildFunc[T]) Module[T] {
	var t T
	return Module[T]{
		name:    contextKey(reflect.TypeOf(t).String()),
		buildFn: buildFunc,
	}
}

func (m Module[T]) Value(ctx context.Context) (ret T) {
	v := ctx.Value(m)
	if v == nil {
		return
	}

	return v.(T)
}

func (m Module[T]) String() string {
	return fmt.Sprintf("Module[%s]", m.name)
}

func (m Module[T]) contextKey() contextKey {
	return m.name
}

func (m Module[T]) build(ctx *buildContext) (err error) {
	t := ctx.Value(m.name)
	if t != nil {
		return nil
	}

	bctx := ctx.Child(m)

	defer func() {
		p := recover()
		if p == nil {
			return
		}

		e, ok := p.(errBuildError)
		if !ok {
			panic(p)
		}

		err = fmt.Errorf("Module[%s] build error: %w", e.name, e.err)
	}()

	instance, err := m.buildFn(bctx)
	if err != nil {
		return fmt.Errorf("Module[%s] build error: %w", m.name, err)
	}

	bctx.addInstance(instance)

	return nil
}
