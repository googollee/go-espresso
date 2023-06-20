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
	Debug(msg string, args ...any)
	DebugCtx(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
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

type EngineOption func(s *Engine) error

func WithCodec(defaultCodec Codec, codec ...Codec) EngineOption {
	return func(s *Engine) error {
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

func WithLogger(logger Logger) EngineOption {
	return func(s *Engine) error {
		if logger == nil {
			err := errors.New("logger should not be nil")
			panic(err)
		}
		s.logger = logger
		return nil
	}
}

type Engine struct {
	httpServer   *http.Server
	router       *httprouter.Router
	defaultCodec Codec
	codecs       map[string]Codec
	logger       Logger
}

func NewEngine(options ...EngineOption) (*Engine, error) {
	ret := &Engine{
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

func (s *Engine) ListenAndServe(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Engine) WithPrefix(prefix string) Router {
	return &router{
		svr:    s,
		prefix: "/" + strings.Trim(prefix, "/"),
	}
}

func (s *Engine) server() *Engine {
	return s
}

func (s *Engine) handle(method, path string, fn httprouter.Handle) {
	s.router.Handle(method, path, fn)
}

func (s *Engine) codec(mime string) Codec {
	if s.codecs != nil {
		if codec, ok := s.codecs[mime]; ok {
			return codec
		}
	}

	return s.defaultCodec
}

type router struct {
	svr    *Engine
	prefix string
}

func (r *router) WithPrefix(prefix string) Router {
	return &router{
		svr:    r.svr,
		prefix: r.prefix + "/" + strings.Trim(prefix, "/"),
	}
}

func (r *router) handle(method, path string, fn httprouter.Handle) {
	path = r.prefix + "/" + strings.TrimLeft(path, "/")
	r.svr.router.Handle(method, path, fn)
}

func (r *router) server() *Engine {
	return r.svr
}

type Router interface {
	WithPrefix(prefix string) Router
	server() *Engine
	handle(method, path string, fn httprouter.Handle)
}
