package espresso

import (
	"context"
	"net/http"
)

type Handler[ContextData any] func(Context[ContextData])

type EndpointDeclarator interface {
	BindPathParam(name string, v any) EndpointDeclarator
	// BindQueryParam(name string, v any) EndpointDeclarator
	// BindHeader(name string, v any) EndpointDeclarator
	End()
}

type Context[Data any] interface {
	context.Context
	Data() *Data
	Endpoint(method, path string, middlewares ...Handler[Data]) EndpointDeclarator

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Abort(err error)
	Error() error
	Next()
}
