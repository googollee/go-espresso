package espresso

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type Handler[Data any] func(Context[Data]) error

func generateHandler[Data any](server *Engine, ctx *declareContext[Data], init Data, fn Handler[Data]) httprouter.Handle {
	endpoint := ctx.endpoint
	ctx.brew.handlers = append(ctx.brew.handlers, fn)

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer r.Body.Close()

		brew := ctx.brew
		ctx := brewContext[Data]{
			Context:  r.Context(),
			brewing:  &brew,
			logger:   server.logger,
			endpoint: endpoint,
			request:  r,
			responserWriter: &responseWriter[Data]{
				ResponseWriter: w,
			},
			pathParams: params,
			data:       init,
		}
		ctx.responserWriter.ctx = &ctx
		brew.ctx = &ctx

		mime := endpoint.ResponseMime
		if mime == "" {
			mime = server.defaultCodec.Mime()
		}
		if mime != "" {
			ctx.responserWriter.Header().Set("Content-Type", mime)
		}

		ctx.Next()

		if ctx.hasWroteResponseCode {
			return
		}

		if ctx.error == nil {
			ctx.responserWriter.WriteHeader(http.StatusNoContent)
			return
		}

		code := http.StatusInternalServerError
		err := ctx.error
		var coder HTTPCoder
		if ok := errors.As(err, &coder); ok {
			code = coder.HTTPCode()
		}

		ctx.responserWriter.WriteHeader(code)
		server.defaultCodec.NewEncoder(ctx.responserWriter).Encode(err)
	}
}

func Handle[Data any](r Router, init Data, fn Handler[Data]) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	declareContext := &declareContext[Data]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		_ = fn(declareContext)
	}()

	endpoint := declareContext.endpoint
	r.server().logger.Info("espresso handles", "method", endpoint.Method, "path", endpoint.Path)
	r.handle(endpoint.Method, endpoint.Path, generateHandler(r.server(), declareContext, init, fn))
}

func HandleProcedure[Data, Request, Response any](r Router, init Data, fn func(Context[Data], *Request) (*Response, error)) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	var req Request

	declareContext := &declareContext[Data]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		_, _ = fn(declareContext, &req)
	}()

	endpoint := declareContext.endpoint
	r.handle(endpoint.Method, endpoint.Path, generateHandler(r.server(), declareContext, init, func(ctx Context[Data]) error {
		mime := ctx.Request().Header.Get("Content-Type")
		codec := r.server().codec(mime)

		var req Request
		if err := codec.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
			return ErrWithStatus(http.StatusBadRequest, err)
		}

		resp, err := fn(ctx, &req)
		if err != nil {
			return err
		}

		if err := codec.NewEncoder(ctx.ResponseWriter()).Encode(resp); err != nil {
			return ErrWithStatus(http.StatusInternalServerError, err)
		}

		return nil
	}))
}

func HandleConsumer[Data, Request any](r Router, init Data, fn func(Context[Data], *Request) error) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	var req Request

	declareContext := &declareContext[Data]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		_ = fn(declareContext, &req)
	}()

	endpoint := declareContext.endpoint
	r.handle(endpoint.Method, endpoint.Path, generateHandler(r.server(), declareContext, init, func(ctx Context[Data]) error {
		mime := ctx.Request().Header.Get("Content-Type")
		codec := r.server().codec(mime)

		var req Request
		if err := codec.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
			return ErrWithStatus(http.StatusBadRequest, err)
		}

		err := fn(ctx, &req)
		if err != nil {
			return err
		}

		return nil
	}))
}

func HandleProvider[Data, Response any](r Router, init Data, fn func(Context[Data]) (*Response, error)) {
	t := reflect.TypeOf(init)
	if t.Kind() == reflect.Ptr {
		panic("ContextData must NOT be a reference type, nor a pointer.")
	}

	declareContext := &declareContext[Data]{
		Context: context.Background(),
	}

	func() {
		defer func() {
			r := recover()
			if _, ok := r.(declareChcecker); ok {
				return
			}

			panic(r) // repanic other values.
		}()
		_, _ = fn(declareContext)
	}()

	endpoint := declareContext.endpoint
	r.handle(endpoint.Method, endpoint.Path, generateHandler(r.server(), declareContext, init, func(ctx Context[Data]) error {
		mime := ctx.Request().Header.Get("Content-Type")
		codec := r.server().codec(mime)

		resp, err := fn(ctx)
		if err != nil {
			return err
		}

		if err := codec.NewEncoder(ctx.ResponseWriter()).Encode(resp); err != nil {
			return ErrWithStatus(http.StatusInternalServerError, err)
		}

		return nil
	}))
}
