package module

import (
	"context"
	"fmt"
	"reflect"

	"github.com/googollee/go-espresso"
)

type ModuleImplementer interface {
	CheckHealthy(context.Context) error
}

type ModuleBuilder[T ModuleImplementer] func(context.Context, *espresso.Server) (T, error)

type ModuleType[T ModuleImplementer] struct {
	name    nameKey
	builder ModuleBuilder[T]
	depends []nameKey
	zeroT   T
}

func NewModule[T ModuleImplementer](builder ModuleBuilder[T]) *ModuleType[T] {
	var t T
	name := reflect.TypeOf(t).String()

	return &ModuleType[T]{
		name:    nameKey(name),
		builder: builder,
	}
}

func (m ModuleType[T]) Name() nameKey {
	return m.name
}

func (m ModuleType[T]) DependOn() []nameKey {
	return m.depends
}

func (m ModuleType[T]) Value(ctx context.Context) T {
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

func (m *ModuleType[T]) CheckHealthy(ctx context.Context) (err error) {
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

func (m *ModuleType[T]) valueWithBuilder(ctx *buildContext) T {
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

func (m *ModuleType[T]) build(ctx *buildContext) error {
	bctx := ctx.Child()

	ret, err := m.builder(bctx, bctx.server)
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
