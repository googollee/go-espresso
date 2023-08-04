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

		reqCodec, respCodec := s.codecs.decideCodec(r)
		respWriter := responseWriter{
			ResponseWriter: w,
			logger:         logger,
		}

		var err error
		ctx := runtimeContext{
			Context:        r.Context(),
			request:        r,
			responseWriter: &respWriter,
			pathParams:     p,
			logger:         logger,
			reqCodec:       reqCodec,
			respCodec:      respCodec,
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

	var failed any

	if runtimeError != nil {
		failed = runtimeError
		Error(ctx, "Error", "error", runtimeError)
	}

	if panicErr != nil {
		failed = panicErr
		Error(ctx, "Panic", "recover", panicErr)
	}

	if w.wroteHeader {
		return
	}

	if failed == nil {
		w.ensureWriteHeader()
		return
	}

	code := http.StatusInternalServerError
	if coder, ok := failed.(HTTPCoder); ok {
		code = coder.HTTPCode()
	}

	codec := ctx.respCodec
	w.Header().Add("Content-Type", codec.Mime())
	w.WriteHeader(code)
	if err := codec.Encode(w, failed); err != nil {
		Error(ctx, "Write response", "error", err)
	}
}
