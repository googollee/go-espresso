package espresso

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var ErrModuleDependError = errors.New("depend error")

type Module interface {
	Name() moduleName
	DependOn() []moduleName
	CheckHealthy(context.Context) error
}

type ModuleImplementer interface {
	CheckHealthy(context.Context) error
}

type moduleName string

type ModuleType[T ModuleImplementer] struct {
	name    moduleName
	builder func(context.Context, *Server) (T, error)
	depends []moduleName
}

func NewModule[T ModuleImplementer](builder func(context.Context, *Server) (T, error)) *ModuleType[T] {
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

type buildContext struct {
	context.Context
	deps map[moduleName]struct{}
}

func (c *buildContext) Value(key any) any {
	name, ok := key.(moduleName)
	if !ok {
		return c.Context.Value(key)
	}

	c.deps[name] = struct{}{}
	return c.Context.Value(name)
}

func (m *ModuleType[T]) Build(ctx context.Context) (T, error) {
	buildContext := &buildContext{
		Context: ctx,
		deps:    make(map[moduleName]struct{}),
	}
	ret, err := m.builder(buildContext)
	if err != nil {
		return ret, err
	}

	m.depends = make([]moduleName, 0, len(buildContext.deps))
	for name := range buildContext.deps {
		m.depends = append(m.depends, name)
	}

	return ret, nil
}

func (m ModuleType[T]) Value(ctx context.Context) T {
	var n T
	v := ctx.Value(m.name)
	if v == nil {
		return n
	}

	ret, ok := v.(T)
	if !ok {
		return n
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

func CheckHealthy(ctx context.Context, rootModules []Module) map[moduleName]error {
	ret := map[moduleName]error{}

	checkModuleHealthy(ctx, rootModules, ret)

	return ret
}

func checkModuleHealthy(ctx context.Context, modules []Module, errs map[moduleName]error) {
	for _, m := range modules {
		if _, ok := errs[m.Name()]; ok {
			continue
		}

		depHealthy := true
		if deps := m.DependOn(); len(deps) != 0 {
			checkModuleHealthy(ctx, deps, errs)

			for _, dep := range deps {
				if errs[dep.Name()] != nil {
					depHealthy = false
					break
				}
			}
		}

		if !depHealthy {
			errs[m.Name()] = ErrModuleDependError
			continue
		}

		errs[m.Name()] = m.CheckHealthy(ctx)
	}
}
