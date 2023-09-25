package espresso

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var ErrModuleDependError = errors.New("depend error")
var ErrModuleNotFound = errors.New("not found module")

type moduleName string

type Module interface {
	Name() moduleName
	DependOn() []moduleName
	CheckHealthy(context.Context) error

	build(*buildContext) error
}

type ModuleImplementer interface {
	CheckHealthy(context.Context) error
}

type ModuleBuilder[T ModuleImplementer] func(context.Context, *Server) (T, error)

type ModuleType[T ModuleImplementer] struct {
	name    moduleName
	builder ModuleBuilder[T]
	depends []moduleName
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

func CheckHealthy(ctx context.Context, rootNames []moduleName) map[moduleName]error {
	ret := map[moduleName]error{}

	checkModuleHealthy(ctx, rootNames, ret)

	return ret
}

func checkModuleHealthy(ctx context.Context, names []moduleName, errs map[moduleName]error) {
	for _, name := range names {
		if _, ok := errs[name]; ok {
			continue
		}

		v := ctx.Value(name)
		if v == nil {
			errs[name] = fmt.Errorf("module %s: %w", name, ErrModuleNotFound)
			continue
		}
		module, ok := v.(Module)
		if !ok {
			errs[name] = fmt.Errorf("module %s: %w", name, ErrModuleNotFound)
			continue
		}

		depHealthy := true
		if deps := module.DependOn(); len(deps) != 0 {
			checkModuleHealthy(ctx, deps, errs)

			for _, depname := range deps {
				if errs[depname] != nil {
					depHealthy = false
					break
				}
			}
		}

		if !depHealthy {
			errs[name] = ErrModuleDependError
			continue
		}

		errs[name] = module.CheckHealthy(ctx)
	}
}

type buildContext struct {
	context.Context
	server  *Server
	deps    map[moduleName]struct{}
	modules map[moduleName]ModuleImplementer
}

func (c *buildContext) Value(key any) any {
	name, ok := key.(moduleName)
	if !ok {
		return c.Context.Value(key)
	}

	c.deps[name] = struct{}{}
	return c.modules[name]
}

func (m *ModuleType[T]) build(ctx *buildContext) error {
	if _, ok := ctx.modules[m.Name()]; ok {
		return nil
	}

	ret, err := m.builder(ctx, ctx.server)
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
