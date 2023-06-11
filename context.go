package espresso

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type EndpointDeclarator interface {
	BindPathParam(name string, v any) EndpointDeclarator
	// BindQueryParam(name string, v any) EndpointDeclarator
	// BindHeader(name string, v any) EndpointDeclarator
	End()
}

type Context[Data any] interface {
	context.Context
	WithContext(ctx context.Context) Context[Data]
	Logger() Logger
	WithLogger(logger Logger) Context[Data]
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
	endpoint *endpoint
	brew     brew[Data]
}

func (c *declareContext[Data]) panic() {
	panic("ctx.Endpoint().BindXXX().End() should be called in the beginning, which is not.")
}

func (c *declareContext[Data]) WithContext(ctx context.Context) Context[Data] {
	c.panic()
	return nil
}

func (c *declareContext[Data]) WithLogger(logger Logger) Context[Data] {
	c.panic()
	return nil
}

func (c *declareContext[Data]) Logger() Logger {
	c.panic()
	return nil
}

func (c *declareContext[Data]) Request() *http.Request {
	c.panic()
	return nil
}

func (c *declareContext[Data]) ResponseWriter() http.ResponseWriter {
	c.panic()
	return nil
}

func (c *declareContext[Data]) Data() *Data {
	c.panic()
	return nil
}

func (c *declareContext[Data]) Abort() {
	c.panic()
}

func (c *declareContext[Data]) Error() error {
	c.panic()
	return nil
}

func (c *declareContext[Data]) Next() {
	c.panic()
}

type brewContext[Data any] struct {
	context.Context
	brewing
	logger          Logger
	endpoint        *endpoint
	request         *http.Request
	responserWriter *responseWriter[Data]
	pathParams      httprouter.Params
	data            Data

	hasWroteResponseCode bool
	isAborted            bool
	error                error
}

func (c *brewContext[Data]) WithContext(ctx context.Context) Context[Data] {
	ret := *c
	ret.Context = ctx
	return &ret
}

func (c *brewContext[Data]) WithLogger(logger Logger) Context[Data] {
	ret := *c
	ret.logger = logger
	return &ret
}

func (c *brewContext[Data]) Logger() Logger {
	return c.logger
}

func (c *brewContext[Data]) Request() *http.Request {
	return c.request
}

func (c *brewContext[Data]) ResponseWriter() http.ResponseWriter {
	return c.responserWriter
}

func (c *brewContext[Data]) Data() *Data {
	return &c.data
}

func (c *brewContext[Data]) Abort() {
	c.isAborted = true
}

func (c *brewContext[Data]) Error() error {
	return c.error
}
