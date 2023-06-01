package espresso

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type processContext[Data any] struct {
	context.Context
	data               *Data
	request            *http.Request
	responseWriter     http.ResponseWriter
	pathParams         httprouter.Params
	endpointDeclarator *endpointDeclarator[Data]
	handlers           []Handler[Data]

	isAborted bool
	err       error
}

func newProcessContext[Data any](ctx context.Context, d *endpointDeclarator[Data], req *http.Request, writer http.ResponseWriter, params httprouter.Params, initData Data, handlers []Handler[Data]) Context[Data] {
	return &processContext[Data]{
		Context:            ctx,
		data:               &initData,
		request:            req,
		responseWriter:     writer,
		pathParams:         params,
		endpointDeclarator: d,
		handlers:           handlers,
	}
}

func (c *processContext[Data]) BindURLParam(name string, v any) error {
	return nil
}

func (c *processContext[Data]) BindQueryParam(name string, v any) error {
	return nil
}

func (c *processContext[Data]) BindHead(name string, v any) error {
	return nil
}

func (c *processContext[Data]) Endpoint(method, path string, middlewares ...Handler[Data]) EndpointDeclarator {
	return newContextBind[Data](c)
}

func (c *processContext[Data]) Data() *Data {
	return c.data
}

func (c *processContext[Data]) Abort(err error) {
	c.isAborted = true
	c.err = err
}

func (c *processContext[Data]) Error() error {
	return c.err
}

func (c *processContext[Data]) Next() {
	for len(c.handlers) > 0 && !c.isAborted {
		handler := c.handlers[0]
		c.handlers = c.handlers[1:]
		handler(c)
	}
}

func (c *processContext[Data]) Request() *http.Request {
	return c.request
}

func (c *processContext[Data]) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}
