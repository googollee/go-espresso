package espresso

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type Logger interface {
	With(args ...any) Logger
	WithGroup(name string) Logger
	DebugCtx(ctx context.Context, msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	ErrorCtx(ctx context.Context, msg string, args ...any)
}

type defaultLogger struct {
	*slog.Logger
}

func (l defaultLogger) With(args ...any) Logger {
	return defaultLogger{
		Logger: l.Logger.With(args...),
	}
}

func (l defaultLogger) WithGroup(name string) Logger {
	return defaultLogger{
		Logger: l.Logger.WithGroup(name),
	}
}

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

func WithLogger(logger Logger) ServerOption {
	return func(s *Server) error {
		if logger == nil {
			err := errors.New("logger should not be nil")
			panic(err)
		}
		s.logger = logger
		return nil
	}
}

type Server struct {
	httpServer   *http.Server
	router       *httprouter.Router
	defaultCodec Codec
	codecs       map[string]Codec
	logger       Logger
}

func NewServer(options ...ServerOption) (*Server, error) {
	ret := &Server{
		httpServer:   &http.Server{},
		router:       httprouter.New(),
		defaultCodec: CodecJSON,
		logger:       defaultLogger{Logger: slog.Default()},
	}

	ret.httpServer.Handler = ret.router

	for _, opt := range options {
		if err := opt(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (s *Server) ListenAndServe(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) WithPrefix(prefix string) Router {
	return &router{
		svr:    s,
		prefix: strings.TrimRight(prefix, "/"),
	}
}

func (s *Server) Handle(method, path string, fn httprouter.Handle) {
	s.router.Handle(method, path, fn)
}

func (s *Server) server() *Server {
	return s
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
	svr    *Server
	prefix string
}

func (r *router) Handle(method, path string, fn httprouter.Handle) {
	path = r.prefix + "/" + strings.TrimLeft(path, "/")
	r.svr.router.Handle(method, path, fn)
}

func (r *router) server() *Server {
	return r.svr
}

type Router interface {
	Handle(method, path string, fn httprouter.Handle)
	server() *Server
}
