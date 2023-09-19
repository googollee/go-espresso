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
	CheckHealthy(context.Context) error
}

type ModuleImplementer interface {
	CheckHealthy(context.Context) error
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
