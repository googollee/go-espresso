package espresso

import (
	"context"
	"fmt"
	"reflect"
)

type moduleName string

type Module[T any] struct {
	name moduleName
}

func DefineModule[T any]() Module[T] {
	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		panic("T should be a type, not a pointer.")
	}

	name := fmt.Sprintf("%T", t)

	return Module[T]{
		name: moduleName(name),
	}
}

func (m Module[T]) With(ctx context.Context, moduleInstance *T) context.Context {
	return context.WithValue(ctx, m.name, moduleInstance)
}

func (m Module[T]) Value(ctx context.Context) *T {
	v := ctx.Value(m.name)
	if v == nil {
		return nil
	}

	ret, ok := v.(*T)
	if !ok {
		return nil
	}

	return ret
}
