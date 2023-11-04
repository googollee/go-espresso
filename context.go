package espresso

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/googollee/go-espresso/module"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const (
	contextRequestCodec  contextKey = "request_codec"
	contextResponseCodec contextKey = "response_codec"
	contextLogger        contextKey = "logger"
)

func getValueFromContext[T any](ctx context.Context, key any) T {
	var zero T

	v := ctx.Value(key)
	if v == nil {
		return zero
	}

	ret, ok := v.(T)
	if !ok {
		return zero
	}

	return ret

}

func RequestCodec(ctx context.Context) Codec {
	return getValueFromContext[Codec](ctx, contextRequestCodec)
}

func ResponseCodec(ctx context.Context) Codec {
	return getValueFromContext[Codec](ctx, contextResponseCodec)
}

func Logger(ctx context.Context) *slog.Logger {
	return getValueFromContext[*slog.Logger](ctx, contextLogger)
}

type Context interface {
	context.Context
	Endpoint(method, path string) EndpointBuilder

	Logger() *slog.Logger
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Abort()
	Error() error
	Next()
}

func WithContext(ctx Context, new context.Context) Context {
	ret := getRuntimeContext(ctx)
	ret.Context = new
	return ret
}

func WithResponseWriter(ctx Context, w http.ResponseWriter) Context {
	ret := getRuntimeContext(ctx)
	ret.responseWriter = w
	return ret
}

func WithLogAttr(ctx Context, args ...any) Context {
	ret := getRuntimeContext(ctx)
	ret.logger = ret.logger.With(args...)
	return ret
}

type runtimeContext struct {
	context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
	pathParams     httprouter.Params
	repo           *module.Repo

	logger    *slog.Logger
	reqCodec  Codec
	respCodec Codec
	endpoint  *Endpoint
	abort     bool
	handlers  []HandleFunc
	err       *error
}

func getRuntimeContext(ctx Context) *runtimeContext {
	rCtx, ok := ctx.(*runtimeContext)
	if !ok {
		panic(errRegisterContextCall)
	}

	return rCtx
}

func (c *runtimeContext) Value(key any) any {
	moduleKey, ok := key.(module.ModuleKey)
	if ok {
		if ret := c.repo.Value(moduleKey); ret != nil {
			return ret
		}
	}

	switch key {
	case contextRequestCodec:
		return c.request
	case contextResponseCodec:
		return c.respCodec
	case contextLogger:
		return c.logger
	}

	return c.Context.Value(key)
}

func (c *runtimeContext) Endpoint(method, path string) EndpointBuilder {
	return &endpointBinder{
		context:  c,
		endpoint: c.endpoint,
	}
}

func (c *runtimeContext) Logger() *slog.Logger {
	return c.logger
}

func (c *runtimeContext) Request() *http.Request {
	return c.request
}

func (c *runtimeContext) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *runtimeContext) Abort() {
	c.abort = true
}

func (c *runtimeContext) Error() error {
	return *c.err
}

func (c *runtimeContext) Next() {
	for !c.abort && len(c.handlers) > 0 {
		handler := c.handlers[0]
		c.handlers = c.handlers[1:]

		if err := handler(c); err != nil {
			*c.err = err
			c.abort = true
		}
	}
}

var errRegisterContextCall = errors.New("call Context.Endpoint() in the beginning with calling End().")

type registerContext struct {
	context.Context
	endpoint *Endpoint
}

func (c *registerContext) Endpoint(method, path string) EndpointBuilder {
	c.endpoint.Method = method
	c.endpoint.Path = path

	return &endpointBuilder{
		endpoint: c.endpoint,
	}
}

func (c *registerContext) Logger() *slog.Logger {
	panic(errRegisterContextCall)
}

func (c *registerContext) Request() *http.Request {
	panic(errRegisterContextCall)
}

func (c *registerContext) ResponseWriter() http.ResponseWriter {
	panic(errRegisterContextCall)
}

func (c *registerContext) Abort() {
	panic(errRegisterContextCall)
}

func (c *registerContext) Error() error {
	panic(errRegisterContextCall)
}

func (c *registerContext) Next() {
	panic(errRegisterContextCall)
}
