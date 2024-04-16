package espresso

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func RPC[Request, Response any](fn func(Context, Request) (Response, error)) HandleFunc {
	return func(ctx Context) error {
		var req Request
		if bctx, ok := ctx.(*buildtimeContext); ok {
			bctx.endpoint.RequestType = reflect.TypeOf(&req).Elem()
			var resp Response
			bctx.endpoint.ResponseType = reflect.TypeOf(&resp).Elem()

			_, err := fn(bctx, req)
			return err
		}

		codec := CodecsModule.Value(ctx)
		if codec == nil {
			return Error(http.StatusInternalServerError, errors.New("no codec in the context"))
		}

		if err := codec.DecodeRequest(ctx, &req); err != nil {
			return Error(http.StatusBadRequest, fmt.Errorf("can't decode request: %w", err))
		}

		resp, err := fn(ctx, req)
		if err != nil {
			return err
		}

		if err := codec.EncodeResponse(ctx, &resp); err != nil {
			return Error(http.StatusInternalServerError, fmt.Errorf("can't encode response: %w", err))
		}

		return nil
	}
}

func RPCRetrive[Response any](fn func(Context) (Response, error)) HandleFunc {
	return func(ctx Context) error {
		if bctx, ok := ctx.(*buildtimeContext); ok {
			var resp Response
			bctx.endpoint.ResponseType = reflect.TypeOf(&resp).Elem()

			_, err := fn(bctx)
			return err
		}

		codec := CodecsModule.Value(ctx)
		if codec == nil {
			return Error(http.StatusInternalServerError, errors.New("no codec in the context"))
		}

		resp, err := fn(ctx)
		if err != nil {
			return err
		}

		if err := codec.EncodeResponse(ctx, &resp); err != nil {
			return Error(http.StatusInternalServerError, fmt.Errorf("can't encode response: %w", err))
		}

		return nil
	}
}

func RPCConsume[Request any](fn func(Context, Request) error) HandleFunc {
	return func(ctx Context) error {
		var req Request
		if bctx, ok := ctx.(*buildtimeContext); ok {
			bctx.endpoint.RequestType = reflect.TypeOf(&req).Elem()

			err := fn(bctx, req)
			return err
		}

		codec := CodecsModule.Value(ctx)
		if codec == nil {
			return Error(http.StatusInternalServerError, errors.New("no codec in the context"))
		}

		if err := codec.DecodeRequest(ctx, &req); err != nil {
			return Error(http.StatusBadRequest, fmt.Errorf("can't decode request: %w", err))
		}

		err := fn(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}
}
