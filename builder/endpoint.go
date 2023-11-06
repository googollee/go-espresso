package builder

import (
	"errors"

	"github.com/googollee/go-espresso/basetype"
)

var ErrEndpointBuildEnd = errors.New("endpoint built")

type endpointBuilder struct {
	endpoint *basetype.Endpoint
	err      basetype.BindErrors
}

func (e *endpointBuilder) BindPath(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bind := e.bindParam(key, basetype.BindPathParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.PathParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) BindForm(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bind := e.bindParam(key, basetype.BindFormParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.FormParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) BindQuery(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bind := e.bindParam(key, basetype.BindQueryParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.QueryParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) BindHead(key string, v any, opts ...basetype.BindOption) basetype.EndpointBuilder {
	bind := e.bindParam(key, basetype.BindHeadParam, v, opts)
	if bind != nil {
		e.addBindParam(&e.endpoint.HeadParams, key, *bind)
	}
	return e
}

func (e *endpointBuilder) End() basetype.BindErrors {
	if len(e.err) != 0 {
		panic(e.err)
	}

	panic(ErrEndpointBuildEnd)
}

func (e *endpointBuilder) bindParam(key string, src basetype.BindSource, v any, opts []basetype.BindOption) *basetype.BindParam {
	bind, err := newBindParam(key, src, v)
	if err != nil {
		e.err = append(e.err, basetype.ErrBind(bind, err))
		return nil
	}

	for _, opt := range opts {
		if err := opt(&bind); err != nil {
			e.err = append(e.err, basetype.ErrBind(bind, err))
			return nil
		}
	}

	return &bind
}

func (e *endpointBuilder) addBindParam(m *map[string]basetype.BindParam, key string, b basetype.BindParam) {
	if *m == nil {
		*m = make(map[string]basetype.BindParam)
	}
	(*m)[key] = b
}
