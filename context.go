package espresso

import (
	"context"
	"net/http"
)

type Handler[ContextData any] func(Context[ContextData])

type Context[ContextData any] interface {
	context.Context
	Data() *ContextData
	BindURLParam(name string, v any) error
	BindQueryParam(name string, v any) error
	BindHead(name string, v any) error
	Endpoint(method, path string, middlewares ...Handler[ContextData])

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Abort(err error)
	Error() error
	Next()
}
