package codec

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type Codec interface {
	Mime() string
	EncodeResponse(ctx Context, v any) error
	DecodeRequest(ctx Context, v any) error
}
