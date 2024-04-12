package espresso

import (
	"fmt"
	"net/http"

	"github.com/googollee/go-espresso/codec"
)

func cacheAllError(ctx Context) error {
	wr := &responseWriter{
		ResponseWriter: ctx.ResponseWriter(),
	}
	code := http.StatusInternalServerError
	defer func() {
		perr := recover()

		if wr.hasWritten || (ctx.Error() == nil && perr == nil) {
			return
		}

		if httpCoder, ok := ctx.Error().(HTTPError); ok {
			code = httpCoder.HTTPCode()
		}
		wr.WriteHeader(code)

		if perr == nil {
			perr = ctx.Error()
		}

		codec := codec.Module.Value(ctx)
		if codec == nil {
			fmt.Fprintf(wr, "%v", perr)
			return
		}

		_ = codec.EncodeResponse(ctx, perr)
	}()

	ctx = ctx.WithResponseWriter(wr)
	ctx.Next()

	return nil
}
