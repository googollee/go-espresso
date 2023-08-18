package espresso

import (
	"errors"
	"net/http"
)

func Procedure[Response, Request any](ctx Context, fn func(Context, Request) (Response, error)) error {
	var req Request

	rCtx, ok := ctx.(*runtimeContext)
	if !ok {
		_, _ = fn(ctx, req)
		return nil
	}

	err := rCtx.reqCodec.Decode(rCtx.Request().Body, &req)
	if err != nil {
		err = ErrWithStatus(http.StatusBadRequest, err)
	}

	var resp Response
	if err == nil {
		resp, err = fn(rCtx, req)
	}

	rCtx.responseWriter.Header().Add("Content-Type", rCtx.respCodec.Mime())

	code := http.StatusOK
	var ret any = resp
	if err != nil {
		ret = err

		code = http.StatusInternalServerError
		var hc HTTPCoder
		if errors.As(err, &hc) {
			code = hc.HTTPCode()
		}
	}
	rCtx.ResponseWriter().WriteHeader(code)

	if err := rCtx.respCodec.Encode(rCtx.ResponseWriter(), ret); err != nil {
		return err
	}

	return nil
}

func Consume[Request any](ctx Context, fn func(Context, Request) error) error {
	var req Request

	rCtx, ok := ctx.(*runtimeContext)
	if !ok {
		_ = fn(ctx, req)
		return nil
	}

	if err := rCtx.reqCodec.Decode(rCtx.Request().Body, &req); err != nil {
		return ErrWithStatus(http.StatusBadRequest, err)
	}

	if err := fn(rCtx, req); err != nil {
		return err
	}

	rCtx.ResponseWriter().WriteHeader(http.StatusNoContent)
	return nil
}

func Provide[Response any](ctx Context, fn func(Context) (Response, error)) error {
	rCtx, ok := ctx.(*runtimeContext)
	if !ok {
		_, _ = fn(ctx)
		return nil
	}

	resp, err := fn(rCtx)
	if err != nil {
		return err
	}

	rCtx.responseWriter.Header().Add("Content-Type", rCtx.respCodec.Mime())
	rCtx.ResponseWriter().WriteHeader(http.StatusOK)
	if err := rCtx.respCodec.Encode(rCtx.ResponseWriter(), resp); err != nil {
		return err
	}

	return nil

}
