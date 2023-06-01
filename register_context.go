package espresso

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type registerContext[Data any] struct {
	context.Context
}

func newRegisterContext[Data any](router *httprouter.Router) Context[Data] {
	return &registerContext[Data]{}
}

func (c *registerContext[Data]) BindURLParam(name string, v any) error {
	return nil
}

func (c *registerContext[Data]) BindQueryParam(name string, v any) error {
	return nil
}

func (c *registerContext[Data]) BindHead(name string, v any) error {
	return nil
}

func (c *registerContext[Data]) Endpoint(method, path string, middlewares ...Handler[Data]) EndpointDeclarator {
	return newEndpointDeclarator[Data](method, path, middlewares)
}

func (c *registerContext[Data]) Data() *Data {
	return nil
}

func (c *registerContext[Data]) Abort(err error) {
}

func (c *registerContext[Data]) Error() error {
	return nil
}

func (c *registerContext[Data]) Next() {
}

func (c *registerContext[Data]) Request() *http.Request {
	return nil
}

func (c *registerContext[Data]) ResponseWriter() http.ResponseWriter {
	return nil
}
