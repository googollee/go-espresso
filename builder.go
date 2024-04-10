package espresso

import (
	"context"
	"errors"
	"net/http"
)

var errBuilderEnd = errors.New("build end.")
var errRegisterContextCall = errors.New("should call Context.Endpoint() in the beginning with End().")

type buildEndpoint struct {
	endpoint *Endpoint
	err      BindErrors
}

func (b *buildEndpoint) BindPath(key string, v any) EndpointBuilder {
	bind, err := newBindParam(key, BindPathParam, v)
	if err != nil {
		panic(errorBind(bind, err).Error())
	}

	b.endpoint.PathParams[key] = bind

	return b
}

func (b *buildEndpoint) End() BindErrors {
	panic(errBuilderEnd)
}

type buildContext struct {
	context.Context
	endpoint *Endpoint
}

func newBuildContext() Context {
	return &buildContext{
		Context:  context.Background(),
		endpoint: &Endpoint{},
	}
}

func (c *buildContext) Endpoint(method, path string, middlewares ...HandleFunc) EndpointBuilder {
	c.endpoint.Method = method
	c.endpoint.Path = path
	c.endpoint.MiddlewareFuncs = middlewares

	return &buildEndpoint{
		endpoint: c.endpoint,
	}
}

func (c *buildContext) Request() *http.Request {
	panic(errRegisterContextCall)
}

func (c *buildContext) ResponseWriter() http.ResponseWriter {
	panic(errRegisterContextCall)
}

func (c *buildContext) Error() error {
	panic(errRegisterContextCall)
}

func (c *buildContext) Next() {
	panic(errRegisterContextCall)
}

func (c *buildContext) WithParent(ctx context.Context) Context {
	panic(errRegisterContextCall)
}

func (c *buildContext) WithResponseWriter(w http.ResponseWriter) Context {
	panic(errRegisterContextCall)
}
