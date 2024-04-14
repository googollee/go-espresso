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
		err := checkError(ctx, recover())

		if wr.hasWritten || err == nil {
			return
		}

		if httpCoder, ok := err.(HTTPError); ok {
			code = httpCoder.HTTPCode()
		}
		wr.WriteHeader(code)

		codec := codec.Module.Value(ctx)
		if codec == nil {
			fmt.Fprintf(wr, "%v", err)
			return
		}

		_ = codec.EncodeResponse(ctx, err)
	}()

	ctx = ctx.WithResponseWriter(wr)
	ctx.Next()

	return nil
}

func checkError(ctx Context, perr any) error {
	if perr != nil {
		return Error(http.StatusInternalServerError, fmt.Errorf("%v", perr))
	}

	err := ctx.Err()
	if err == nil {
		return nil
	}

	if _, ok := err.(HTTPError); ok {
		return err
	}

	return Error(http.StatusInternalServerError, err)
}
