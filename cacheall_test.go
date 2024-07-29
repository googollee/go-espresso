package espresso_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/googollee/module"

	"github.com/googollee/go-espresso"
)

func TestCacheAllMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		providers   []module.Provider
		middlewares []espresso.HandleFunc
		wantCode    int
		wantBody    string
	}{
		{
			name: "MiddlewareError",
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				return errors.New("error")
			}},
			wantCode: http.StatusInternalServerError,
			wantBody: "error",
		},
		{
			name: "MiddlewareHTTPError",
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				return espresso.Error(http.StatusGatewayTimeout, errors.New("gateway timeout"))
			}},
			wantCode: http.StatusGatewayTimeout,
			wantBody: "gateway timeout",
		},
		{
			name: "MiddlewarePanic",
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				panic("panic")
			}},
			wantCode: http.StatusInternalServerError,
			wantBody: "panic",
		},
		{
			name:      "MiddlewareErrorWithCodec",
			providers: []module.Provider{espresso.ProvideCodecs},
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				return errors.New("error")
			}},
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"message\":\"error\"}\n",
		},
		{
			name:      "MiddlewareHTTPErrorWithCodec",
			providers: []module.Provider{espresso.ProvideCodecs},
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				return espresso.Error(http.StatusGatewayTimeout, errors.New("gateway timeout"))
			}},
			wantCode: http.StatusGatewayTimeout,
			wantBody: "{\"message\":\"gateway timeout\"}\n",
		},
		{
			name:      "MiddlewarePanicWithCodec",
			providers: []module.Provider{espresso.ProvideCodecs},
			middlewares: []espresso.HandleFunc{func(ctx espresso.Context) error {
				panic("panic")
			}},
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"message\":\"panic\"}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			espo := espresso.New()
			espo.Use(tc.middlewares...)
			espo.AddModule(tc.providers...)

			var called int32
			espo.HandleFunc(func(ctx espresso.Context) error {
				atomic.AddInt32(&called, 1)

				if err := ctx.Endpoint(http.MethodGet, "/").End(); err != nil {
					return err
				}

				fmt.Fprint(ctx.ResponseWriter(), "ok")

				return nil
			})

			called = 0
			svr := httptest.NewServer(espo)
			defer svr.Close()

			resp, err := http.Get(svr.URL)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if got := atomic.LoadInt32(&called); got != 0 {
				t.Fatalf("handle func is called")
			}

			if got, want := resp.StatusCode, tc.wantCode; got != want {
				t.Fatalf("resp.Status = %d, want: %d", got, want)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := string(respBody), tc.wantBody; got != want {
				t.Errorf("resp.Body = %q, want: %q", got, want)
			}
		})
	}
}
