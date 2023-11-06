package builder

import (
	"context"
	"errors"
	"net/http"

	"github.com/googollee/go-espresso/basetype"
)

var errRegisterContextCall = errors.New("should call Context.Endpoint() in the beginning with End().")

type Context struct {
	context.Context
	EndpointDef basetype.Endpoint
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

func (c *Context) Endpoint(method, path string, fn ...basetype.HandleFunc) basetype.EndpointBuilder {
	c.EndpointDef.Method = method
	c.EndpointDef.Path = path
	c.EndpointDef.MiddlewareFuncs = fn

	return &endpointBuilder{
		endpoint: &c.EndpointDef,
	}
}

func (c *Context) Request() *http.Request {
	panic(errRegisterContextCall)
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return nilResponsWriter{}
}

func (c *Context) Error() error {
	panic(errRegisterContextCall)
}

func (c *Context) Next() {
	panic(errRegisterContextCall)
}

func (c *Context) WithParent(ctx context.Context) basetype.Context {
	panic(errRegisterContextCall)
}

func (c *Context) WithResponseWriter(w http.ResponseWriter) basetype.Context {
	panic(errRegisterContextCall)
}

type nilResponsWriter struct{}

func (w nilResponsWriter) Header() http.Header {
	return http.Header{}
}

func (w nilResponsWriter) Write(p []byte) (int, error) {
	panic(errRegisterContextCall)
}

func (w nilResponsWriter) WriteHeader(code int) {
	panic(errRegisterContextCall)
}
