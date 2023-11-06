package builder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextPanic(t *testing.T) {
	ctx := NewContext(context.Background())

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
		if got, want := v, errRegisterContextCall; got != want {
			t.Errorf("call %s got panic %q, want: %q", name, got, want)
		}
	}

	func() {
		defer recover("Request()")
		ctx.Request()
	}()
	func() {
		defer recover("ResponseWriter().WriteHeader()")
		ctx.ResponseWriter().WriteHeader(http.StatusOK)
	}()
	func() {
		defer recover("ResponseWriter().Write()")
		_, _ = ctx.ResponseWriter().Write([]byte("hello"))
	}()
	func() {
		defer recover("Error()")
		_ = ctx.Error()
	}()
	func() {
		defer recover("Next()")
		ctx.Next()
	}()
	func() {
		defer recover("WithParent()")
		ctx.WithParent(context.Background())
	}()
	func() {
		defer recover("WithResponseWriter()")
		ctx.WithResponseWriter(httptest.NewRecorder())
	}()
}

func TestContextResponseWriterSetHeader(t *testing.T) {
	ctx := NewContext(context.Background())

	ctx.ResponseWriter().Header().Set("Content-Type", "application/json")
}
