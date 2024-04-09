package espresso

import (
	"context"
	"net/http"
)

type HandleFunc func(Context) error

type ContextExtender interface {
	WithParent(ctx context.Context) Context
	WithResponseWriter(w http.ResponseWriter) Context
}

type Context interface {
	context.Context
	ContextExtender
	Endpoint(method string, path string, middlewares ...HandleFunc) EndpointBuilder

	Error() error
	Next()

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type MiddlewareProvider interface {
	Middlewares() []HandleFunc
}
