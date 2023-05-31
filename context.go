package espresso

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler[ContextData any] func(*Context[ContextData])

type Context[Data any] struct {
	context.Context
	Data Data

	request        *http.Request
	responseWriter http.ResponseWriter
	pathParams     httprouter.Params

	handlers  []Handler[Data]
	err       error
	isAborted bool
}

func newContext[Data any](ctx context.Context, req *http.Request, writer http.ResponseWriter, params httprouter.Params, initData Data, handlers []Handler[Data]) *Context[Data] {
	return &Context[Data]{
		Context:        ctx,
		Data:           initData,
		request:        req,
		responseWriter: writer,
		pathParams:     params,
		handlers:       handlers,
	}
}

func (c *Context[Data]) Abort(err error) {
	c.isAborted = true
	c.err = err
}

func (c *Context[Data]) Error() error {
	return c.err
}

func (c *Context[Data]) Next() {
	for len(c.handlers) > 0 && !c.isAborted {
		handler := c.handlers[0]
		c.handlers = c.handlers[1:]
		handler(c)
	}
}

func (c *Context[Data]) Request() *http.Request {
	return c.request
}

func (c *Context[Data]) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context[Data]) PathParam(name string) string {
	return c.pathParams.ByName(name)
}
