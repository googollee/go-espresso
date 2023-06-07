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

type ServerOption func(s *Server) error

func WithCodec(defaultCodec Codec, codec ...Codec) ServerOption {
	return func(s *Server) error {
		if s.codecs == nil {
			s.codecs = make(map[string]Codec)
		}

		s.defaultCodec = defaultCodec
		s.codecs[defaultCodec.Mime()] = defaultCodec

		for _, c := range codec {
			s.codecs[c.Mime()] = c
		}

		return nil
	}
}

type Server struct {
	server       *http.Server
	router       *httprouter.Router
	defaultCodec Codec
	codecs       map[string]Codec
}

func NewServer(options ...ServerOption) (*Server, error) {
	ret := &Server{
		server:       &http.Server{},
		router:       httprouter.New(),
		defaultCodec: CodecJSON,
	}

	ret.server.Handler = ret.router

	for _, opt := range options {
		if err := opt(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
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

func (s *Server) codec(mime string) Codec {
	if s.codecs != nil {
		if codec, ok := s.codecs[mime]; ok {
			return codec
		}
	}

	return s.defaultCodec
}

type router struct {
	server *Server
	prefix string
}

type Handler[Data any] func(Context[Data]) error

func (r *router) Handle(method, path string, fn httprouter.Handle) {
	path = r.prefix + "/" + strings.TrimLeft(path, "/")
	r.server.router.Handle(method, path, fn)
}

func (r *router) codec(mime string) Codec {
	return r.server.codec(mime)
}

type Router interface {
	Handle(method, path string, fn httprouter.Handle)
	codec(mime string) Codec
}

func generateHandler[Data any](ctx *declareContext[Data], init Data, fn Handler[Data]) httprouter.Handle {
	endpoint := ctx.endpoint
	ctx.brew.handlers = append(ctx.brew.handlers, fn)

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer r.Body.Close()

		brew := ctx.brew
		ctx := brewContext[Data]{
			Context:  r.Context(),
			Brewing:  &brew,
			endpoint: endpoint,
			request:  r,
			responserWriter: &responseWriter[Data]{
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
	}
}

func Handle[Data any](r Router, init Data, fn Handler[Data]) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	declareContext := &declareContext[Data]{
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
	fmt.Println("handle", endpoint.Method, endpoint.Path)
	r.Handle(endpoint.Method, endpoint.Path, generateHandler(declareContext, init, fn))
}

func HandleProcedure[Data, Request, Response any](r Router, init Data, fn func(Context[Data], *Request) (*Response, error)) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	var req Request

	declareContext := &declareContext[Data]{
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
		_, _ = fn(declareContext, &req)
	}()

	endpoint := declareContext.endpoint
	fmt.Println("handle", endpoint.Method, endpoint.Path)
	r.Handle(endpoint.Method, endpoint.Path, generateHandler(declareContext, init, func(ctx Context[Data]) error {
		codec := r.codec("")

		var req Request
		if err := codec.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
			return WithStatus(http.StatusBadRequest, err)
		}

		resp, err := fn(ctx, &req)
		if err != nil {
			return err
		}

		if err := codec.NewEncoder(ctx.ResponseWriter()).Encode(resp); err != nil {
			return WithStatus(http.StatusInternalServerError, err)
		}

		return nil
	}))
}

func HandleConsumer[Data, Request any](r Router, init Data, fn func(Context[Data], Request) error) {
}

func HandleProvider[Data, Response any](r Router, init Data, fn func(Context[Data]) (Response, error)) {
}
