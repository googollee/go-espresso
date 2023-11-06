package runtime

import (
	"context"
	"errors"
	"net/http"

	"github.com/googollee/go-espresso/basetype"
	"github.com/googollee/go-espresso/module"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Repo     *module.Repo
	Endpoint basetype.Endpoint
	Handlers []basetype.HandleFunc
	Error    *error
}

type Context struct {
	context.Context
	server Server

	request        *http.Request
	responseWriter http.ResponseWriter
	pathParams     httprouter.Params
	abort          bool
}

func NewContext(ctx context.Context, s Server, w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		Context:        ctx,
		server:         s,
		request:        r,
		responseWriter: w,
		pathParams:     params,
	}
}

func (c *Context) Value(key any) any {
	moduleKey, ok := key.(module.Key)
	if ok {
		if ret := c.server.Repo.Value(c, moduleKey); ret != nil {
			return ret
		}
	}

	return c.Context.Value(key)
}

func (c *Context) Endpoint(method, path string, fn ...basetype.HandleFunc) basetype.EndpointBuilder {
	return &EndpointBinder{
		context:  c,
		endpoint: c.server.Endpoint,
	}
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) Error() error {
	return *c.server.Error
}

func (c *Context) Next() {
	if c.abort || len(c.server.Handlers) == 0 {
		return
	}

	handler := c.server.Handlers[0]
	c.server.Handlers = c.server.Handlers[1:]

	if err := handler(c); err != nil {
		c.abort = true
		if *c.server.Error == nil {
			*c.server.Error = err
		} else {
			*c.server.Error = errors.Join(*c.server.Error, err)
		}
	}
}

func (c *Context) WithParent(parent context.Context) basetype.Context {
	return &Context{
		Context:        parent,
		server:         c.server,
		request:        c.request,
		responseWriter: c.responseWriter,
		pathParams:     c.pathParams,
		abort:          c.abort,
	}
}

func (c *Context) WithResponseWriter(w http.ResponseWriter) basetype.Context {
	return &Context{
		Context:        c.Context,
		server:         c.server,
		request:        c.request,
		responseWriter: w,
		pathParams:     c.pathParams,
		abort:          c.abort,
	}
}
