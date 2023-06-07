package espresso

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	server *http.Server
	router *httprouter.Router
}

func NewServer() *Server {
	ret := &Server{
		server: &http.Server{},
		router: httprouter.New(),
	}

	ret.server.Handler = ret.router

	return ret
}

func (s *Server) ListenAndServe(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) WithPrefix(prefix string) Router {
	return &router{
		server: s,
		prefix: strings.TrimRight(prefix, "/"),
	}
}

func (s *Server) Handle(method, path string, fn httprouter.Handle) {
	s.router.Handle(method, path, fn)
}

type router struct {
	server *Server
	prefix string
}

func (r *router) Handle(method, path string, fn httprouter.Handle) {
	path = r.prefix + "/" + strings.TrimLeft(path, "/")
	r.server.router.Handle(method, path, fn)
}

type Router interface {
	Handle(method, path string, fn httprouter.Handle)
}

func Handle[ContextData any](r Router, init ContextData, fn Handler[ContextData]) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	declareContext := &declareContext[ContextData]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		_ = fn(declareContext)
	}()

	endpoint := declareContext.endpoint
	declareContext.brew.handlers = append(declareContext.brew.handlers, fn)

	fmt.Println("handle", endpoint.Method, endpoint.Path)
	r.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		brew := declareContext.brew
		ctx := brewContext[ContextData]{
			Context:  r.Context(),
			Brewing:  &brew,
			endpoint: endpoint,
			request:  r,
			responserWriter: &responseWriter[ContextData]{
				ResponseWriter: w,
			},
			pathParams: params,
			data:       init,
		}
		ctx.responserWriter.ctx = &ctx
		brew.ctx = &ctx

		ctx.Next()

		if ctx.hasWroteResponseCode {
			return
		}

		if ctx.error == nil {
			ctx.responserWriter.WriteHeader(http.StatusOK)
			return
		}

		code := http.StatusInternalServerError
		var coder HTTPCoder
		if ok := errors.As(ctx.error, &coder); ok {
			code = coder.HTTPCode()
		}
		ctx.responserWriter.WriteHeader(code)
		_, _ = ctx.responserWriter.Write([]byte(ctx.error.Error()))
	})
}
