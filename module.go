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
	Check(context.Context) error
}

type ModuleImplementer interface {
	Check(context.Context) error
}

type moduleName string

type ModuleType[T ModuleImplementer] struct {
	name    moduleName
	depends []Module
}

func DefineModule[T ModuleImplementer](depends ...Module) *ModuleType[T] {
	var t T
	name := reflect.TypeOf(t).Name()
	return &ModuleType[T]{
		name:    moduleName(name),
		depends: depends,
	}
}

func (m ModuleType[T]) Name() moduleName {
	return m.name
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

func (m *ModuleType[T]) Check(ctx context.Context) error {
	var errs []error
	for _, module := range m.depends {
		if err := module.Check(ctx); err != nil {
			errs = append(errs, fmt.Errorf("module %s: %w", module.Name(), err))
		}
	}

	if len(errs) != 0 {
		errs = append(errs, fmt.Errorf("module %s: %w", m.Name(), ErrModuleDependError))
	} else if err := m.Check(ctx); err != nil {
		errs = append(errs, fmt.Errorf("module %s: %w", m.Name(), err))
	}

	return errors.Join(errs...)
}
