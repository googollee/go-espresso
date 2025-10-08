package espresso

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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

func handleRPC(fn reflect.Value) (HandleFunc, error) {
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("fn(%T) should be a function.", fn)
	}
	tfn := fn.Type()

	var inType reflect.Type
	if tfn.NumIn() == 2 {
		inType = tfn.In(1).Elem()
	}
	var outType reflect.Type
	if tfn.NumOut() == 2 {
		outType = tfn.Out(0).Elem()
	}

	return func(ctx Context) error {
		var req reflect.Value
		bctx, isBuildTime := ctx.(*buildtimeContext)

		if inType != nil {
			req = reflect.New(inType)
			if isBuildTime {
				bctx.endpoint.RequestType = req.Elem().Type()
			}
		}

		inputs := []reflect.Value{reflect.ValueOf(ctx)}
		if inType != nil {
			inputs = append(inputs, req)
		}

		if isBuildTime {
			outputs := fn.Call(inputs)
			if outType != nil {
				return outputs[1].Interface().(error)
			}
			return outputs[0].Interface().(error)
		}

		codec := CodecsModule.Value(ctx)

		if inType != nil {
			if codec == nil {
				return Error(http.StatusInternalServerError, errors.New("no codec in the context"))
			}

			if err := codec.DecodeRequest(ctx, req.Interface()); err != nil {
				return Error(http.StatusBadRequest, fmt.Errorf("can't decode request: %w", err))
			}
		}

		outputs := fn.Call(inputs)
		var err error
		if outType != nil {
			if v := outputs[1].Interface(); v != nil {
				err = v.(error)
			}
		} else {
			if v := outputs[0].Interface(); v != nil {
				err = v.(error)
			}
		}

		if err != nil {
			return err
		}

		if outType.Kind() == reflect.Invalid {
			return nil
		}

		resp := outputs[0]
		if err := codec.EncodeResponse(ctx, resp.Interface()); err != nil {
			fmt.Fprintln(os.Stderr, "can't encode")
			return Error(http.StatusInternalServerError, fmt.Errorf("can't encode response: %w", err))
		}

		return nil
	}, nil
}
