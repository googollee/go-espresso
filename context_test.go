package espresso

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/exp/slog"
)

type fakeKey string

var checkKey fakeKey = "key"

type fakeContext struct {
	context.Context
	name string
}

func (c fakeContext) Value(key any) any {
	if key == checkKey {
		return c.name
	}
	return nil
}

func TestWithContext(t *testing.T) {
	req := &http.Request{}
	resp := httptest.NewRecorder()
	logger := slog.Default()
	err := errors.New("error")

	var rctx Context = &runtimeContext{
		Context:        fakeContext{name: "base"},
		request:        req,
		responseWriter: resp,
		logger:         logger,
		err:            &err,
	}

	if got, want := rctx.Value(checkKey), "base"; got != want {
		t.Errorf("rctx.Value(checkKey) = %q, want: %q", got, want)
	}

	rctx = WithContext(rctx, fakeContext{name: "new"})

	if got, want := rctx.Value(checkKey), "new"; got != want {
		t.Errorf("rctx.Value(checkKey) = %q, want: %q", got, want)
	}
	if got, want := rctx.Request(), req; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.ResponseWriter(), resp; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Logger(), logger; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Error(), err; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
}

func TestWithResponseWriter(t *testing.T) {
	req := &http.Request{}
	oldResp := httptest.NewRecorder()
	logger := slog.Default()
	err := errors.New("error")

	var rctx Context = &runtimeContext{
		Context:        fakeContext{name: "base"},
		request:        req,
		responseWriter: oldResp,
		logger:         logger,
		err:            &err,
	}

	if got, want := rctx.ResponseWriter(), oldResp; got != want {
		t.Errorf("rctx.ResponseWriter() = %p, want: %p", got, want)
	}

	newResp := httptest.NewRecorder()
	rctx = WithResponseWriter(rctx, newResp)
	if got, want := rctx.Value(checkKey), "base"; got != want {
		t.Errorf("rctx.Value(checkKey) = %q, want: %q", got, want)
	}
	if got, want := rctx.Request(), req; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.ResponseWriter(), newResp; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Logger(), logger; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Error(), err; got != want {
		t.Errorf("rctx.Request() = %p, want: %p", got, want)
	}
}

func TestWithLogAttr(t *testing.T) {
	var buf bytes.Buffer
	logHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return a
		},
	})

	req := &http.Request{}
	resp := httptest.NewRecorder()
	logger := slog.New(logHandler)
	err := errors.New("error")

	var rctx Context = &runtimeContext{
		Context:        fakeContext{name: "base"},
		request:        req,
		responseWriter: resp,
		logger:         logger,
		err:            &err,
	}

	Info(rctx, "message")
	if got, want := buf.String(), "{\"level\":\"INFO\",\"msg\":\"message\"}\n"; got != want {
		t.Errorf("Info(rctx, \"message\"), log = %q, want: %q", got, want)
	}

	buf.Reset()
	rctx = WithLogAttr(rctx, "spec", "key")
	Info(rctx, "message")
	if got, want := buf.String(), "{\"level\":\"INFO\",\"msg\":\"message\",\"spec\":\"key\"}\n"; got != want {
		t.Errorf("Info(rctx, \"message\"), log = %q, want: %q", got, want)
	}
}

func TestRegisterContextPanic(t *testing.T) {
	var ctx Context = &registerContext{
		endpoint: &Endpoint{},
	}

	func() {
		defer func() {
			if got := recover(); got != nil {
				t.Errorf("call Endpoint() got panic %q, want: nil", got)
			}
		}()
		_ = ctx.Endpoint("", "")
	}()

	recover := func(name string) {
		v := recover()
		if got, want := fmt.Sprintf("%v", v), "call Context.Endpoint() in the beginning with calling End()."; got != want {
			t.Errorf("call %s got panic %q, want: %q", name, got, want)
		}
	}

	func() {
		defer recover("Logger()")
		ctx.Logger()
	}()
	func() {
		defer recover("Request()")
		ctx.Request()
	}()
	func() {
		defer recover("ResponseWriter()")
		_ = ctx.ResponseWriter()
	}()
	func() {
		defer recover("Abort()")
		ctx.Abort()
	}()
	func() {
		defer recover("Error()")
		_ = ctx.Error()
	}()
	func() {
		defer recover("Next()")
		ctx.Next()
	}()
}
