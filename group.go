package espresso

import (
	"fmt"
	"reflect"
	"strings"
)

type Group struct {
	prefix      string
	middlewares []HandleFunc
	server      *Server
}

func (g *Group) WithPrefix(path string) *Group {
	return &Group{
		prefix:      strings.TrimRight(g.prefix, "/") + "/" + strings.Trim(path, "/"),
		middlewares: g.middlewares[0:len(g.middlewares)],
		server:      g.server,
	}
}

func (g *Group) Use(middleware ...HandleFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *Group) HandleAll(svc any) {
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

func (g *Group) HandleFunc(fn HandleFunc) {
	g.handleFunc(fn, fmt.Sprintf("%T", fn))
}

func (g *Group) handleFunc(fn HandleFunc, sig string) {
	var endpoint Endpoint
	ctx := registerContext{
		Context:  nil,
		endpoint: &endpoint,
	}

	defer func() {
		v := recover()
		if v != errEndpointBuildEnd {
			if v == nil {
				v = errRegisterContextCall
			}
			panic(v)
		}
		endpoint.Path = strings.TrimRight(g.prefix, "/") + "/" + strings.Trim(endpoint.Path, "/")

		g.server.registerEndpoint(&endpoint, g.middlewares, fn, sig)
	}()

	_ = fn(&ctx)
}
