package espresso

import "net/http"

func Produce[Response, Request any](ctx Context, fn func(Context, Request) (Response, error)) error {
	var req Request

	if _, ok := ctx.(*registerContext); ok {
		_, _ = fn(ctx, req)
		return nil
	}

	codec := DefaultCodec
	if err := codec.Decode(ctx.Request().Body, &req); err != nil {
		return ErrWithStatus(http.StatusBadRequest, err)
	}

	resp, err := fn(ctx, req)
	if err != nil {
		return err
	}

	ctx.ResponseWriter().WriteHeader(http.StatusOK)
	if err := codec.Encode(ctx.ResponseWriter(), resp); err != nil {
		return err
	}

	return nil
}

func Consume[Request any](ctx Context, fn func(Context, Request) error) error {
	var req Request

	if _, ok := ctx.(*registerContext); ok {
		_ = fn(ctx, req)
		return nil
	}

	codec := DefaultCodec
	if err := codec.Decode(ctx.Request().Body, &req); err != nil {
		return ErrWithStatus(http.StatusBadRequest, err)
	}

	if err := fn(ctx, req); err != nil {
		return err
	}

	ctx.ResponseWriter().WriteHeader(http.StatusNoContent)
	return nil
}

func Provide[Response any](ctx Context, fn func(Context) (Response, error)) error {
	if _, ok := ctx.(*registerContext); ok {
		_, _ = fn(ctx)
		return nil
	}

	codec := DefaultCodec
	resp, err := fn(ctx)
	if err != nil {
		return err
	}

	ctx.ResponseWriter().WriteHeader(http.StatusOK)
	if err := codec.Encode(ctx.ResponseWriter(), resp); err != nil {
		return err
	}

	return nil

}
