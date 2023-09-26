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
	name    moduleName
	builder ModuleBuilder[T]
	depends []moduleName
	zeroT   T
}

func NewModule[T ModuleImplementer](builder ModuleBuilder[T]) *ModuleType[T] {
	var t T
	name := reflect.TypeOf(t).Name()

	return &ModuleType[T]{
		name:    moduleName(name),
		builder: builder,
	}
}

func (m ModuleType[T]) Name() moduleName {
	return m.name
}

func (m ModuleType[T]) DependOn() []moduleName {
	return m.depends
}

func (m ModuleType[T]) Value(ctx context.Context) T {
	v := ctx.Value(m.name)
	if v == nil {
		return m.zeroT
	}

	if bctx, ok := ctx.(*buildContext); ok {
		return m.valueWithBuilder(bctx)
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
	if ret, ok := ctx.modules[m.Name()]; ok {
		return ret.(T)
	}

	if err := m.build(ctx); err != nil {
		ctx.err = fmt.Errorf("module %s: %q", m.Name(), err)
		panic(errBuildError)
	}

	return ctx.modules[m.Name()].(T)
}

func (m *ModuleType[T]) build(ctx *buildContext) error {
	bctx := buildContext{
		Context: ctx.Context,
		server:  ctx.server,
		deps:    make(map[moduleName]struct{}),
		modules: ctx.modules,
		err:     nil,
	}

	ret, err := m.builder(&bctx, bctx.server)
	if err != nil {
		return err
	}

	m.depends = make([]moduleName, 0, len(ctx.deps))
	for name := range ctx.deps {
		m.depends = append(m.depends, name)
	}

	ctx.modules[m.Name()] = ret

	return nil
}
