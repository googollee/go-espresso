package espresso

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func WithoutDefault() ServerOption {
	return func(s *Server) error {
		s.middlewares = nil
		return nil
	}
}

type responseWriter struct {
	http.ResponseWriter
	responseCode     int
	hasWrittenHeader bool
}

func (w *responseWriter) WriteHeader(code int) {
	w.responseCode = code
	w.hasWrittenHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.hasWrittenHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(b)
}

func MiddlewareLogRequest(ctx Context) error {
	req := ctx.Request()
	ctx = WithLogAttr(ctx, "span", time.Now().Unix(), "method", req.Method, "path", req.URL.Path)

	respW := &responseWriter{
		ResponseWriter: ctx.ResponseWriter(),
	}
	ctx = WithResponseWriter(ctx, respW)

	Info(ctx, "Request")
	defer func() {
		var ret any
		code := http.StatusInternalServerError

		if err := ctx.Error(); err != nil {
			Error(ctx, "Error", "error", err.Error())

			ret = err

			var hc HTTPCoder
			if errors.As(err, &hc) {
				code = hc.HTTPCode()
			}
		}
		if p := recover(); p != nil {
			Error(ctx, "Panic", "panic", fmt.Sprintf("%s", p))

			ret = p
			code = http.StatusInternalServerError
		}

		if !respW.hasWrittenHeader {
			codec := ResponseCodec(ctx)
			if codec == nil {
				return
			}

			respW.Header().Set("Content-Type", codec.Mime())
			respW.WriteHeader(code)
			if err := codec.Encode(respW, ret); err != nil {
				Error(ctx, "Encode response", "error", err.Error())
			}
		}

		Info(ctx, "Response", "code", respW.responseCode)
	}()

	ctx.Next()

	return nil
}
