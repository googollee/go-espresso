package espresso

import (
	"errors"
	"fmt"
)

type EndpointBuilder interface {
	BindPath(key string, v any, opts ...BindOption) EndpointBuilder
	BindForm(key string, v any, opts ...BindOption) EndpointBuilder
	BindQuery(key string, v any, opts ...BindOption) EndpointBuilder
	End() BindErrors
}

type Endpoint struct {
	Method      string
	Path        string
	PathParams  map[string]BindParam
	QueryParams map[string]BindParam
	FormParams  map[string]BindParam
}

type endpointBinder struct {
	context  *runtimeContext
	endpoint *Endpoint
	err      BindErrors
}

func (e *endpointBinder) BindPath(key string, v any, opts ...BindOption) EndpointBuilder {
	bindParam := e.endpoint.PathParams[key]
	str := e.context.pathParams.ByName(key)

	if err := bindParam.fn(e.context, v, str); err != nil {
		e.err = append(e.err, newBindError(bindParam, err))
		return e
	}

	return e
}

func (e *endpointBinder) BindForm(key string, v any, opts ...BindOption) EndpointBuilder {
	bindParam := e.endpoint.FormParams[key]
	req := e.context.request

	if err := req.ParseForm(); err != nil {
		e.err = append(e.err, newBindError(bindParam, err))
		return e
	}

	str := req.FormValue(key)
	if err := bindParam.fn(e.context, v, str); err != nil {
		e.err = append(e.err, newBindError(bindParam, err))
		return e
	}

	return e
}

func (e *endpointBinder) BindQuery(key string, v any, opts ...BindOption) EndpointBuilder {
	bindParam := e.endpoint.QueryParams[key]
	req := e.context.request

	str := req.URL.Query().Get(key)
	if err := bindParam.fn(e.context, v, str); err != nil {
		e.err = append(e.err, newBindError(bindParam, err))
		return e
	}

	return e
}

func (e *endpointBinder) End() BindErrors {
	if len(e.err) != 0 {
		return e.err
	}
	return nil
}

type endpointBuilderFail []error

var errEndpointBuildEnd = errors.New("endpoint built")

func (e endpointBuilderFail) String() string {
	return fmt.Sprintf("%v", []error(e))
}

type endpointBuilder struct {
	endpoint *Endpoint
	err      endpointBuilderFail
}

func (e *endpointBuilder) BindPath(key string, v any, opts ...BindOption) EndpointBuilder {
	bind := e.bindParam(key, BindPathParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.PathParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) BindForm(key string, v any, opts ...BindOption) EndpointBuilder {
	bind := e.bindParam(key, BindFormParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.FormParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) BindQuery(key string, v any, opts ...BindOption) EndpointBuilder {
	bind := e.bindParam(key, BindQueryParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.QueryParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) End() BindErrors {
	if len(e.err) != 0 {
		panic(e.err)
	}

	panic(errEndpointBuildEnd)
}

func (e *endpointBuilder) bindParam(key string, typ BindType, v any, opts []BindOption) *BindParam {
	bind, err := newBindParam(key, typ, v)
	if err != nil {
		e.err = append(e.err, err)
		return nil
	}

	for _, opt := range opts {
		if err := opt(&bind); err != nil {
			e.err = append(e.err, err)
			return nil
		}
	}

	return &bind
}

func (e *endpointBuilder) addBindParam(m *map[string]BindParam, key string, b BindParam) {
	if *m == nil {
		*m = make(map[string]BindParam)
	}
	(*m)[key] = b
}
