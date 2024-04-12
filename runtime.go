package espresso

import (
	"context"
	"net/http"
	"time"
)

type runtimeEndpoint struct {
	request  *http.Request
	endpoint *Endpoint
	err      BindErrors
}

func (e *runtimeEndpoint) BindPath(key string, v any) EndpointBuilder {
	binder, ok := e.endpoint.PathParams[key]
	if !ok {
		return e
	}

	strV := e.request.PathValue(key)
	if err := binder.Func(v, strV); err != nil {
		e.err = append(e.err, errorBind(binder, err))
	}

	return e
}

func (e *runtimeEndpoint) End() BindErrors {
	return e.err
}

type runtimeContext struct {
	ctx      context.Context
	endpoint *Endpoint
	request  *http.Request
	response http.ResponseWriter

	err        error
	chainIndex int
}

func (c *runtimeContext) Endpoint(method, path string, mid ...HandleFunc) EndpointBuilder {
	return &runtimeEndpoint{
		request:  c.request,
		endpoint: c.endpoint,
	}
}

func (c *runtimeContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *runtimeContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *runtimeContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *runtimeContext) Err() error {
	return c.ctx.Err()
}

func (c *runtimeContext) WithParent(ctx context.Context) Context {
	return &runtimeContext{
		ctx:        ctx,
		endpoint:   c.endpoint,
		request:    c.request,
		response:   c.response,
		err:        c.err,
		chainIndex: c.chainIndex,
	}
}

func (c *runtimeContext) WithResponseWriter(w http.ResponseWriter) Context {
	return &runtimeContext{
		ctx:        c.ctx,
		endpoint:   c.endpoint,
		request:    c.request,
		response:   w,
		err:        c.err,
		chainIndex: c.chainIndex,
	}
}

func (c *runtimeContext) Request() *http.Request {
	return c.request
}

func (c *runtimeContext) ResponseWriter() http.ResponseWriter {
	return c.response
}

func (c *runtimeContext) Next() {
	index := c.chainIndex
	c.chainIndex++
	if err := c.endpoint.ChainFuncs[index](c); err != nil {
		c.err = err
	}
}

func (c *runtimeContext) Error() error {
	return c.err
}
