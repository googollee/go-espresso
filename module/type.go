package module

import (
	"context"
	"fmt"
	"reflect"
)

type Instance interface {
	CheckHealthy(context.Context) error
}

type Builder[T Instance] func(context.Context) (T, error)

type Type[T Instance] struct {
	name    nameKey
	builder Builder[T]
	depends []nameKey
	zeroT   T
}

func NewModule[T Instance](builder Builder[T]) *Type[T] {
	var t T
	name := reflect.TypeOf(t).String()

	return &Type[T]{
		name:    nameKey(name),
		builder: builder,
	}
}

func (m Type[T]) Name() nameKey {
	return m.name
}

func (m Type[T]) DependOn() []nameKey {
	return m.depends
}

func (m Type[T]) Value(ctx context.Context) T {
	if bctx, ok := ctx.(*buildContext); ok {
		return m.valueWithBuilder(bctx)
	}

	v := ctx.Value(m.Name())
	if v == nil {
		return m.zeroT
	}

	ret, ok := v.(T)
	if !ok {
		return m.zeroT
	}

	return ret
}

func (m *Type[T]) CheckHealthy(ctx context.Context) (err error) {
	defer func() {
		v := recover()
		if v == nil {
			return
		}

		err = fmt.Errorf("check healthy panic: %v", v)
	}()

	err = m.Value(ctx).CheckHealthy(ctx)
	return
}

func (m *Type[T]) valueWithBuilder(ctx *buildContext) T {
	ctx.deps[m.Name()] = struct{}{}
	if ret, ok := ctx.instances[m.Name()]; ok {
		return ret.(T)
	}

	if err := m.build(ctx); err != nil {
		ctx.err = fmt.Errorf("module %s: %q", m.Name(), err)
		panic(errBuildError)
	}

	return ctx.instances[m.Name()].(T)
}

func (m *Type[T]) build(ctx *buildContext) error {
	bctx := ctx.Child()

	ret, err := m.builder(bctx)
	if err != nil {
		return err
	}

	m.depends = make([]nameKey, 0, len(ctx.deps))
	for name := range ctx.deps {
		m.depends = append(m.depends, name)
	}

	ctx.instances[m.Name()] = ret

	return nil
}
