package espresso

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
)

type Router interface {
	Use(middlewares ...HandleFunc)
	WithPrefix(path string) Router
	HandleFunc(handleFunc HandleFunc)
	Handle(controller any)
}

type router struct {
	prefix      string
	middlewares []HandleFunc
	mux         *http.ServeMux
}

func (g *router) WithPrefix(path string) Router {
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
	g.handleFunc(fn)
}

func (g *router) handleFunc(fn HandleFunc) {
	ctx := newBuildtimeContext()

	defer func() {
		v := recover()
		if v != errBuilderEnd {
			if v == nil {
				v = fmt.Errorf("should call ctx.Endpoint().End()")
			}
			panic(v)
		}

		g.register(ctx, fn)
	}()

	_ = fn(ctx)
}

func (g *router) register(ctx *buildtimeContext, fn HandleFunc) {
	path := strings.TrimRight(g.prefix, "/") + "/" + strings.TrimLeft(ctx.endpoint.Path, "/")
	chains := slices.Clone(g.middlewares)
	chains = append(chains, ctx.endpoint.ChainFuncs...)
	chains = append(chains, fn)

	endpoint := *ctx.endpoint
	endpoint.Path = path
	endpoint.ChainFuncs = chains

	pattern := ctx.endpoint.Method + " " + path
	g.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &runtimeContext{
			ctx:      r.Context(),
			endpoint: &endpoint,
			request:  r,
			response: w,
		}

		ctx.Next()
	})
}

func (g *router) Handle(controller any) {
	v := reflect.ValueOf(controller)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)

		fn, err := handleRPC(method)
		if err != nil {
			continue
		}

		g.handleFunc(fn)
	}
}
