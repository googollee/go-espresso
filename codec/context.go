package codec

import (
	"context"
	"io"
	"net/http"
)

type Context interface {
	context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type Codec interface {
	Mime() string
	Encode(ctx context.Context, w io.Writer, v any) error
	Decode(ctx context.Context, r io.Reader, v any) error
}
