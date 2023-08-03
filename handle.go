package espresso

import "net/http"

func Produce[Response, Request any](ctx Context, fn func(Context, Request) (Response, error)) error {
	var req Request

	rCtx, ok := ctx.(*runtimeContext)
	if !ok {
		_, _ = fn(ctx, req)
		return nil
	}

	codec := rCtx.codec
	if err := codec.Decode(rCtx.Request().Body, &req); err != nil {
		return ErrWithStatus(http.StatusBadRequest, err)
	}

	resp, err := fn(rCtx, req)
	if err != nil {
		return err
	}

	rCtx.responseWriter.Header().Add("Content-Type", codec.Mime())
	rCtx.ResponseWriter().WriteHeader(http.StatusOK)
	if err := codec.Encode(rCtx.ResponseWriter(), resp); err != nil {
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

	codec := rCtx.codec
	if err := codec.Decode(rCtx.Request().Body, &req); err != nil {
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

	codec := rCtx.codec
	resp, err := fn(rCtx)
	if err != nil {
		return err
	}

	rCtx.responseWriter.Header().Add("Content-Type", codec.Mime())
	rCtx.ResponseWriter().WriteHeader(http.StatusOK)
	if err := codec.Encode(rCtx.ResponseWriter(), resp); err != nil {
		return err
	}

	return nil

}
