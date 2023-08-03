package espresso

import (
	"net/http"
	"time"

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
		logger := s.logger.With(
			"span", time.Now().Unix(),
			"method", r.Method,
			"path", r.URL.Path,
		)

		var err error
		respWriter := responseWriter{
			ResponseWriter: w,
			logger:         logger,
		}

		ctx := runtimeContext{
			Context:        r.Context(),
			request:        r,
			responseWriter: &respWriter,
			pathParams:     p,
			logger:         logger,
			codec:          s.codecs.decideCodec(r),
			endpoint:       endpoint,
			handlers:       handlers,
			err:            &err,
		}

		defer func() {
			s.done(&ctx, &respWriter, recover(), err)
		}()

		Info(&ctx, "Receive")

		ctx.Next()
	})
}

func (s *Server) done(ctx *runtimeContext, w *responseWriter, panicErr any, runtimeError error) {
	defer w.logCode(ctx)

	msg := "Panic"
	fail := panicErr
	if fail == nil {
		msg = "Error"
		fail = runtimeError
	}

	if fail != nil {
		Error(ctx, msg, "error", fail)
	}

	if w.wroteHeader {
		return
	}

	if fail != nil {
		code := http.StatusInternalServerError
		if coder, ok := fail.(HTTPCoder); ok {
			code = coder.HTTPCode()
		}

		w.WriteHeader(code)
		if err := ctx.codec.Encode(w, fail); err != nil {
			Error(ctx, "Write response", "error", err)
		}
	}

	w.ensureWriteHeader()
}
