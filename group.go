package espresso

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
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
		method := v.Method(i).Interface()
		fn, ok := method.(func(Context) error)
		if !ok {
			fmt.Printf("ignore %T\n", method)
			continue
		}

		g.HandleFunc(fn)
	}
}

func (g *Group) HandleFunc(fn HandleFunc) {
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

		g.registerHandle(&endpoint, fn)
	}()

	_ = fn(&ctx)
}

func (g *Group) registerHandle(endpoint *Endpoint, fn HandleFunc) {
	g.server.endpoints = append(g.server.endpoints, *endpoint)

	Info(context.Background(), "Register handle func", "method", endpoint.Method, "path", endpoint.Path, "handle", fmt.Sprintf("%T", fn))
	g.server.router.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := runtimeContext{
			Context:        r.Context(),
			request:        r,
			responseWriter: w,
			pathParams:     p,
			injectedValues: map[InjectKey]any{},
			endpoint:       endpoint,
			abort:          false,
			handlers:       append(g.middlewares[0:len(g.middlewares)], fn),
			err:            nil,
		}
		ctx.Next()

		if err := ctx.Err(); err != nil {
			code := http.StatusInternalServerError
			if coder, ok := err.(HTTPCoder); ok {
				code = coder.HTTPCode()
			}

			w.WriteHeader(code)
			_ = DefaultCodec.Encode(w, err)
		}
	})
}
