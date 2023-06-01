package espresso

import (
	"context"
	"net/http"
)

type Handler[ContextData any] func(Context[ContextData])

type EndpointDeclarator interface {
	BindPathParam(name string, v any) EndpointDeclarator
	// BindQueryParam(name string, v any) EndpointDeclarator
	// BindHead(name string, v any) EndpointDeclarator
}

type Context[ContextData any] interface {
	context.Context
	Data() *ContextData
	Endpoint(method, path string, middlewares ...Handler[ContextData]) EndpointDeclarator

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Abort(err error)
	Error() error
	Next()
}
