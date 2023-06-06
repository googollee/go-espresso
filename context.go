package espresso

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler[ContextData any] func(Context[ContextData]) error

type EndpointDeclarator interface {
	BindPathParam(name string, v any) EndpointDeclarator
	// BindQueryParam(name string, v any) EndpointDeclarator
	// BindHeader(name string, v any) EndpointDeclarator
	End()
}

type Context[Data any] interface {
	context.Context
	Endpoint(method, path string, middlewares ...Handler[Data]) Declarator

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Data() *Data
	Abort()
	Error() error
	Next()
}

type declareContext[Data any] struct {
	context.Context
	endpoint *endpoint[Data]
}

func (c *declareContext[Data]) Request() *http.Request {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) ResponseWriter() http.ResponseWriter {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) Data() *Data {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) Abort() {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) Error() error {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) Next() {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

type handleContext[Data any] struct {
	context.Context
	endpoint        *endpoint[Data]
	handleIndex     int
	request         *http.Request
	responserWriter *responseWriter[Data]
	pathParams      httprouter.Params
	data            Data

	hasWroteResponseCode bool
	isAborted            bool
	error                error
}

func (c *handleContext[Data]) Request() *http.Request {
	return c.request
}

func (c *handleContext[Data]) ResponseWriter() http.ResponseWriter {
	return c.responserWriter
}

func (c *handleContext[Data]) Data() *Data {
	return &c.data
}

func (c *handleContext[Data]) Abort() {
	c.isAborted = true
}

func (c *handleContext[Data]) Error() error {
	return c.error
}

func (c *handleContext[Data]) Next() {
	for c.handleIndex < len(c.endpoint.Handlers) && !c.isAborted {
		handler := c.endpoint.Handlers[c.handleIndex]
		c.handleIndex++
		if err := handler(c); err != nil {
			var ig HTTPIgnore
			if ok := errors.As(err, &ig); ok && ig.Ignore() {
				continue
			}

			c.isAborted = true
			c.error = err
			break
		}
	}
}
