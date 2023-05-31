package rpc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/googollee/go-espresso"
)

type ServiceOpt[ContextData any] func(*Service[ContextData]) error

func WithDefaultMime[ContextData any](codec Codec) ServiceOpt[ContextData] {
	return func(s *Service[ContextData]) error {
		if codec == nil {
			return errors.New("invalid codec nil")
		}

		s.defaultCodec = codec
		return nil
	}
}

func AppendMiddlewares[ContextData any](handlers ...espresso.Handler[ContextData]) ServiceOpt[ContextData] {
	return func(s *Service[ContextData]) error {
		s.handlers = append(s.handlers, handlers...)
		return nil
	}
}

type Service[ContextData any] struct {
	server       *espresso.Server[ContextData]
	defaultCodec Codec
	prefix       string
	handlers     []espresso.Handler[ContextData]
}

func New[ContextData any](server *espresso.Server[ContextData], prefix string, opts ...ServiceOpt[ContextData]) (*Service[ContextData], error) {
	ret := Service[ContextData]{
		server:       server,
		defaultCodec: CodecJSON{},
		prefix:       strings.TrimRight(prefix, "/"),
	}

	for _, opt := range opts {
		if err := opt(&ret); err != nil {
			return nil, fmt.Errorf("configure service error: %w", err)
		}
	}

	return &ret, nil
}

func (s *Service[ContextData]) GetCodec(ctx *espresso.Context[ContextData]) Codec {
	return s.defaultCodec
}

func (s *Service[ContextData]) LoadRequest(ctx *espresso.Context[ContextData], v any) bool {
	r := ctx.Request()
	if err := s.GetCodec(ctx).NewDecoder(r.Body).Decode(v); err != nil {
		s.ResponseError(ctx, err)
		return false
	}

	return true
}

func (s *Service[ContextData]) Response(ctx *espresso.Context[ContextData], code int, v any) {
	w := ctx.ResponseWriter()
	w.Header().Add("Content-ContextDataype", s.defaultCodec.Mime())
	w.WriteHeader(code)

	if v != nil {
		s.GetCodec(ctx).NewEncoder(w).Encode(v)
	}
}

func (s *Service[ContextData]) ResponseError(ctx *espresso.Context[ContextData], err error) {
	code := http.StatusInternalServerError
	if he, ok := err.(HTTPCode); ok {
		code = he.HTTPCode()
	}

	s.Response(ctx, code, err)
}

func (s *Service[ContextData]) WithPath(path string) string {
	return s.prefix + "/" + strings.TrimLeft(path, "/")
}

func GET[ContextData, Response any](svc *Service[ContextData], path string, f func(ctx *espresso.Context[ContextData]) (Response, error)) {
	svc.server.GET(svc.WithPath(path), func(ctx *espresso.Context[ContextData]) {
		resp, err := f(ctx)
		if err != nil {
			svc.ResponseError(ctx, err)
			return
		}

		svc.Response(ctx, http.StatusOK, resp)
	})
}

func POST[ContextData, Request, Response any](svc *Service[ContextData], path string, f func(ctx *espresso.Context[ContextData], req Request) (Response, error)) {
	svc.server.POST(svc.WithPath(path), func(ctx *espresso.Context[ContextData]) {
		var req Request
		if cont := svc.LoadRequest(ctx, &req); !cont {
			return
		}

		resp, err := f(ctx, req)
		if err != nil {
			svc.ResponseError(ctx, err)
			return
		}

		svc.Response(ctx, http.StatusOK, resp)
	})
}
