package espresso

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/googollee/go-espresso/builder"
)

type Router struct {
	prefix      string
	middlewares []HandleFunc
	server      *Server
}

func (g *Router) WithPrefix(path string) *Router {
	return &Router{
		prefix:      strings.TrimRight(g.prefix, "/") + "/" + strings.Trim(path, "/"),
		middlewares: g.middlewares[0:len(g.middlewares)],
		server:      g.server,
	}
}

func (g *Router) Use(middleware ...HandleFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *Router) HandleAll(svc any) {
	v := reflect.ValueOf(svc)

	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		fn, ok := method.Interface().(func(Context) error)
		if !ok {
			continue
		}

		t := v.Type()
		sig := t.Method(i).Name

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		sig = t.String() + "." + sig

		g.handleFunc(fn, sig)
	}
}

func (g *Router) HandleFunc(fn HandleFunc) {
	g.handleFunc(fn, fmt.Sprintf("%T", fn))
}

func (g *Router) handleFunc(fn HandleFunc, sig string) {
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
