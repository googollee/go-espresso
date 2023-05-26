package middleware

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type TestStruct struct {
	Str string
	I   int
}

func TestMiddleware(t *testing.T) {
	m := NewMiddleware[*TestStruct]()
	ctx := &gin.Context{
		Keys: make(map[string]any),
	}

	want := &TestStruct{
		Str: "string",
		I:   10,
	}

	m.Store(ctx, want)
	got := m.Get(ctx)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got != want, diff: (-got, +want):\n%s", diff)
	}
}

func TestMiddlewareKey(t *testing.T) {
	tests := []struct {
		middleware any
		want       string
	}{
		{NewMiddleware[*TestStruct](), "middleware.TestStruct"},
		{NewMiddleware[TestStruct](), "middleware.TestStruct"},
		{NewMiddlewareWithKey[*TestStruct]("mykey"), "mykey"},
	}

	for _, test := range tests {
		if got, want := fmt.Sprint(test.middleware), test.want; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
