package espresso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/googollee/go-espresso/builder"
)

type router struct {
	prefix      string
	middlewares []HandleFunc
	mux         *http.ServeMux
}

func (g *router) WithPrefix(path string) *router {
	return &router{
		prefix:      strings.TrimRight(g.prefix, "/") + "/" + strings.Trim(path, "/"),
		middlewares: g.middlewares[0:len(g.middlewares)],
		mux:         g.mux,
	}
}

func (g *router) Use(middleware ...HandleFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *router) HandleFunc(fn HandleFunc) {
	g.handleFunc(fn, fmt.Sprintf("%T", fn))
}

func (g *router) handleFunc(fn HandleFunc, sig string) {
	var endpoint Endpoint
	ctx := builder.NewContext(context.Background())

	defer func() {
		v := recover()
		if v != builder.ErrEndpointBuildEnd {
			if v == nil {
				v = fmt.Errorf("should call Endpoint().End()")
			}
			panic(v)
		}

		endpoint.Path = strings.TrimRight(g.prefix, "/") + "/" + strings.Trim(endpoint.Path, "/")

		g.server.registerEndpoint(ctx.EndpointDef, g.middlewares, fn, sig)
	}()

	_ = fn(ctx)
}

func (g *router) register(endpoint *Endpoint) {
	path := strings.TrimRight(g.prefix, "/") + "/" + strings.TrimLeft(endpoint.Path, "/")
	middlewares := append(g.middlewares, endpoint.MiddlewareFuncs...)

	g.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
	})
}
