package espresso

import (
	"context"
	"errors"
	"net/http"
)

var errBuilderEnd = errors.New("build end.")
var errRegisterContextCall = errors.New("should call Context.Endpoint() in the beginning with End().")

type buildtimeEndpoint struct {
	endpoint *Endpoint
}

func (b *buildtimeEndpoint) BindPath(key string, v any) EndpointBuilder {
	bind, err := newBindParam(key, BindPathParam, v)
	if err != nil {
		panic(errorBind(bind, err).Error())
	}

	b.endpoint.PathParams[key] = bind

	return b
}

func (b *buildtimeEndpoint) End() BindErrors {
	panic(errBuilderEnd)
}

type buildtimeContext struct {
	context.Context
	endpoint *Endpoint
}

func newBuildtimeContext() *buildtimeContext {
	return &buildtimeContext{
		Context:  context.Background(),
		endpoint: newEndpoint(),
	}
}

func (c *buildtimeContext) Endpoint(method, path string, middlewares ...HandleFunc) EndpointBuilder {
	c.endpoint.Method = method
	c.endpoint.Path = path
	c.endpoint.ChainFuncs = middlewares

	return &buildtimeEndpoint{
		endpoint: c.endpoint,
	}
}

func (c *buildtimeContext) Request() *http.Request {
	panic(errRegisterContextCall)
}

func (c *buildtimeContext) ResponseWriter() http.ResponseWriter {
	panic(errRegisterContextCall)
}

func (c *buildtimeContext) Error() error {
	panic(errRegisterContextCall)
}

func (c *buildtimeContext) Next() {
	panic(errRegisterContextCall)
}

func (c *buildtimeContext) WithParent(ctx context.Context) Context {
	panic(errRegisterContextCall)
}

func (c *buildtimeContext) WithResponseWriter(w http.ResponseWriter) Context {
	panic(errRegisterContextCall)
}
