package espresso

import (
	"context"
	"errors"
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
		t.Errorf("fc.Name() = %q, want: %q", got, want)
	}

	rctx = WithContext(rctx, fakeContext{name: "new"})

	if got, want := rctx.Value(checkKey), "new"; got != want {
		t.Errorf("fc.Name() = %q, want: %q", got, want)
	}
	if got, want := rctx.Request(), req; got != want {
		t.Errorf("fc.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.ResponseWriter(), resp; got != want {
		t.Errorf("fc.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Logger(), logger; got != want {
		t.Errorf("fc.Request() = %p, want: %p", got, want)
	}
	if got, want := rctx.Error(), err; got != want {
		t.Errorf("fc.Request() = %p, want: %p", got, want)
	}
}
