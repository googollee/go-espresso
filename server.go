package espresso

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/googollee/go-espresso/basetype"
	"github.com/googollee/go-espresso/log"
	"github.com/googollee/go-espresso/module"
	"github.com/googollee/go-espresso/runtime"
	"github.com/julienschmidt/httprouter"
)

type BindError = basetype.BindError
type BindErrors = basetype.BindErrors
type Context = basetype.Context
type HandleFunc = basetype.HandleFunc
type Endpoint = basetype.Endpoint
type EndpointBuilder = basetype.EndpointBuilder
type ServerOption = basetype.ServerOption

type Server struct {
	repo *module.Repo

	Router
	endpoints []Endpoint
	router    *httprouter.Router
}

func New(opts ...ServerOption) (*Server, error) {
	ret := &Server{
		router: httprouter.New(),
		repo:   module.NewRepo(),
	}

	ret.Router = Router{
		prefix: "/",
		server: ret,
	}

	for _, opt := range opts {
		if err := opt(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func Default(opts ...ServerOption) (*Server, error) {
	defaultOptions := []ServerOption{
		log.Use(),
	}
	opts = append(defaultOptions, opts...)
	return New(opts...)
}

func (s *Server) AddModule(builders ...module.Builder) {
	s.repo.Add(builders...)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) registerEndpoint(endpoint Endpoint, middle []HandleFunc, fn HandleFunc, fnSignature string) {
	s.endpoints = append(s.endpoints, endpoint)
	handlers := append(middle[0:], fn)

	s.router.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resp := responseWriter{
			ResponseWriter: w,
		}

		var err error
		ctx := runtime.NewContext(r.Context(), runtime.Server{
			Repo:     s.repo,
			Endpoint: endpoint,
			Handlers: handlers,
			Error:    &err,
		}, &resp, r, p)

		ctx.Next()

		if err == nil {
			return
		}

		if resp.hasWritten {
			return
		}

		code := http.StatusInternalServerError
		if hc, ok := err.(HTTPCoder); ok {
			code = hc.HTTPCode()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(err)
	})
}
