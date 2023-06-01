package espresso

import (
	"fmt"

	"github.com/julienschmidt/httprouter"
)

type endpointDeclarator[Data any] struct {
	method   string
	path     string
	handlers []Handler[Data]
	binders  map[string]bindFunc
}

func newEndpointDeclarator[Data any](method, path string, handlers []Handler[Data]) *endpointDeclarator[Data] {
	return &endpointDeclarator[Data]{
		method:   method,
		path:     path,
		handlers: handlers,
		binders:  make(map[string]bindFunc),
	}
}

func (d *endpointDeclarator[Data]) BindPathParam(name string, v any) EndpointDeclarator {
	if _, ok := v.(Binding); ok {
		d.binders[name] = bind
		return d
	}

	switch v.(type) {
	case *int:
		d.binders[name] = bindInt
	case *string:
		d.binders[name] = bindStr
	default:
		panic(fmt.Sprintf("unknown type %T, please try implementing Binding interface.", v))
	}

	return d
}

type contextBind struct {
	params  httprouter.Params
	binders map[string]bindFunc
	errors  map[string]error
}

func newContextBind[Data any](ctx *processContext[Data]) *contextBind {
	return &contextBind{
		params:  ctx.pathParams,
		binders: ctx.endpointDeclarator.binders,
		errors:  make(map[string]error),
	}
}

func (b *contextBind) BindPathParam(name string, v any) EndpointDeclarator {
	bindFunc, ok := b.binders[name]
	if !ok {
		b.errors[name] = fmt.Errorf("no binder for path param with name %q", name)
		return b
	}

	if err := bindFunc(b.params.ByName(name), v); err != nil {
		b.errors[name] = err
		return b
	}

	return b
}
