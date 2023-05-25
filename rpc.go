package espresso

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceOpt func(*Service) error

func WithDefaultMime(codec Codec) ServiceOpt {
	return func(s *Service) error {
		if codec == nil {
			return errors.New("invalid codec nil")
		}

		s.defaultCodec = codec
		return nil
	}
}

type Service struct {
	defaultCodec Codec
}

func NewService(opts ...ServiceOpt) (*Service, error) {
	ret := Service{
		defaultCodec: CodecJSON{},
	}

	for _, opt := range opts {
		if err := opt(&ret); err != nil {
			return nil, fmt.Errorf("configure service error: %w", err)
		}
	}

	return &ret, nil
}

func (s *Service) GetCodec(ctx *gin.Context) Codec {
	return s.defaultCodec
}

func (s *Service) LoadRequest(ctx *gin.Context, v any) bool {
	r := ctx.Request
	if err := s.GetCodec(ctx).NewDecoder(r.Body).Decode(v); err != nil {
		s.ResponseError(ctx, err)
		return false
	}

	return true
}

func (s *Service) Response(ctx *gin.Context, code int, v any) {
	w := ctx.Writer
	w.Header().Add("Content-Type", s.defaultCodec.Mime())
	w.WriteHeader(code)

	if v != nil {
		s.GetCodec(ctx).NewEncoder(w).Encode(v)
	}
}

func (s *Service) ResponseError(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	if he, ok := err.(HTTPCode); ok {
		code = he.HTTPCode()
	}

	s.Response(ctx, code, err)
}

func RPCOnlyResponse[Response any](svc *Service, f func(ctx *gin.Context) (Response, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		resp, err := f(ctx)
		if err != nil {
			svc.ResponseError(ctx, err)
			return
		}

		svc.Response(ctx, http.StatusOK, resp)
	}
}

func RPCOnlyRequest[Request any](svc *Service, f func(ctx *gin.Context, req Request) error) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req Request
		if cont := svc.LoadRequest(ctx, &req); !cont {
			return
		}

		if err := f(ctx, req); err != nil {
			svc.ResponseError(ctx, err)
			return
		}

		svc.Response(ctx, http.StatusNoContent, nil)
	}
}

func RPC[Request, Response any](svc *Service, f func(ctx *gin.Context, req Request) (Response, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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
	}
}
