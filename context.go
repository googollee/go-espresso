package espresso

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Context interface {
	context.Context
	WithContext(ctx context.Context) Context
	Endpoint(method, path string) EndpointBuilder

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Abort()
	Error() error
	Next()
}

type InjectKey string

type Injecting interface {
	InjectValue(key InjectKey, value any)
}

type runtimeContext struct {
	context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
	pathParams     httprouter.Params

	injectedValues map[InjectKey]any
	endpoint       *Endpoint
	abort          bool
	handlers       []HandleFunc
	err            error
}

func (c *runtimeContext) Value(key any) any {
	if str, ok := key.(InjectKey); ok {
		if ret, ok := c.injectedValues[str]; ok {
			return ret
		}
	}

	return c.Context.Value(key)
}

func (c *runtimeContext) InjectValue(key InjectKey, v any) {
	if c.injectedValues == nil {
		c.injectedValues = make(map[InjectKey]any)
	}
	c.injectedValues[key] = v
}

func (c *runtimeContext) Endpoint(method, path string) EndpointBuilder {
	return &endpointBinder{
		context:  c,
		endpoint: c.endpoint,
	}
}

func (c *runtimeContext) WithContext(ctx context.Context) Context {
	ret := *c
	ret.Context = ctx
	return &ret
}

func (c *runtimeContext) Request() *http.Request {
	return c.request
}

func (c *runtimeContext) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *runtimeContext) Abort() {
	c.abort = true
}

func (c *runtimeContext) Error() error {
	return c.err
}

func (c *runtimeContext) Next() {
	for !c.abort && len(c.handlers) > 0 {
		handler := c.handlers[0]
		c.handlers = c.handlers[1:]

		if err := handler(c); err != nil {
			c.err = err
			c.abort = true
		}
	}
}

var errRegisterContextCall = errors.New("call Context.Endpoint() in the beginning with calling End().")

type registerContext struct {
	context.Context
	endpoint *Endpoint
}

func (c *registerContext) Endpoint(method, path string) EndpointBuilder {
	c.endpoint.Method = method
	c.endpoint.Path = path

	return &endpointBuilder{
		endpoint: c.endpoint,
	}
}

func (c *registerContext) WithContext(ctx context.Context) Context {
	panic(errRegisterContextCall)
}

func (c *registerContext) Request() *http.Request {
	panic(errRegisterContextCall)
}

func (c *registerContext) ResponseWriter() http.ResponseWriter {
	panic(errRegisterContextCall)
}

func (c *registerContext) Abort() {
	panic(errRegisterContextCall)
}

func (c *registerContext) Error() error {
	panic(errRegisterContextCall)
}

func (c *registerContext) Next() {
	panic(errRegisterContextCall)
}
