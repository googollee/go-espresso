package espresso

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type HandleFunc func(Context) error

type Server struct {
	logger *slog.Logger

	Group
	endpoints []Endpoint
	router    *httprouter.Router
}

func New() *Server {
	ret := &Server{
		router: httprouter.New(),
		logger: defaultLogger,
	}

	ret.Group = Group{
		prefix: "/",
		server: ret,
	}

	return ret
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	s.logger.Info("Launch espresso server", "addr", addr)
	return http.ListenAndServe(addr, s.router)
}

func funcName(v any) string {
	return runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
}

func (s *Server) registerEndpoint(endpoint *Endpoint, middle []HandleFunc, fn HandleFunc, fnSignature string) {
	s.endpoints = append(s.endpoints, *endpoint)
	handlers := append(middle[0:], fn)

	s.logger.Info("Register handle func", "method", endpoint.Method, "path", endpoint.Path, "handler", fnSignature)

	s.router.Handle(endpoint.Method, endpoint.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger := s.logger.With(
			"span", time.Now().Unix(),
			"method", r.Method,
			"path", r.URL.Path,
		)

		var err error

		ctx := runtimeContext{
			Context:        r.Context(),
			request:        r,
			responseWriter: w,
			pathParams:     p,
			injectedValues: map[InjectKey]any{
				loggerKey: logger,
			},
			endpoint: endpoint,
			handlers: handlers,
			err:      &err,
		}

		defer func() {
			s.done(&ctx, recover(), err)
		}()

		Info(&ctx, "Received request")

		ctx.Next()
	})
}

func (s *Server) done(ctx *runtimeContext, panicErr any, runtimeError error) {
	w := ctx.ResponseWriter()

	msg := "Panic"
	fail := panicErr
	if fail == nil {
		msg = "Error"
		fail = runtimeError
	}
	if fail != nil {
		code := http.StatusInternalServerError
		if err, ok := fail.(error); ok {
			var coder HTTPCoder
			if errors.As(err, &coder) {
				code = coder.HTTPCode()
				fail = coder
			}
		}

		Error(ctx, msg, "code", code, "error", fail)
		w.WriteHeader(code)
		_ = DefaultCodec.Encode(w, fail)
	}

	Info(ctx, "Done")

}
