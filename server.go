package espresso

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type HandleFunc func(Context) error

type ServerOption func(*Server) error

type Server struct {
	logger *slog.Logger
	codecs codecManager

	Group
	endpoints []Endpoint
	router    *httprouter.Router
}

func New(opts ...ServerOption) (*Server, error) {
	ret := &Server{
		router: httprouter.New(),
		logger: defaultLogger,
		codecs: defaultManager(),
	}

	ret.Group = Group{
		prefix: "/",
		server: ret,
	}

	ret.Use(MiddlewareLogRequest)

	for _, opt := range opts {
		if err := opt(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	s.logger.Info("Launch espresso server", "addr", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) registerEndpoint(endpoint *Endpoint, middle []HandleFunc, fn HandleFunc, fnSignature string) {
	s.endpoints = append(s.endpoints, *endpoint)
	handlers := append(middle[0:], fn)

	s.logger.Info("Register", "method", endpoint.Method, "path", endpoint.Path, "handler", fnSignature)

	s.router.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reqCodec, respCodec := s.codecs.decideCodec(r)
		var err error

		ctx := runtimeContext{
			Context:        r.Context(),
			request:        r,
			responseWriter: w,
			pathParams:     p,
			logger:         s.logger,
			reqCodec:       reqCodec,
			respCodec:      respCodec,
			endpoint:       endpoint,
			handlers:       handlers,
			err:            &err,
		}

		ctx.Next()
	})
}
