// An Overall Example
// This example shows basic usage of `espresso` framework. It provides
// endpoints to access `Blog` web, as well as APIs to handle `Blog` data. For
// simplicity and focus, `Blog` stores in memory, with a map instance, and all
// code are in one package.
package espresso_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/basetype"
	"github.com/googollee/go-espresso/module"
)

type logModule struct {
	log    *slog.Logger
	spanID atomic.Int64
}

func (l *logModule) CheckHealth(ctx context.Context) error {
	return nil
}

var LogModule = module.New[*logModule]()

type logKeyType string

var logKey = logKeyType("custom_logger")

func withLogModule(s basetype.Server) error {
	mod := &logModule{
		log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove time from the output for predictable test output.
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		})),
	}

	s.AddModule(LogModule.Builder(func(context.Context) (*logModule, error) { return mod, nil }))

	s.Use(func(ctx basetype.Context) error {
		req := ctx.Request()

		mod := LogModule.Value(ctx)
		ctxMod := &logModule{
			log: mod.log.With("method", req.Method, "path", req.URL.Path, "span", mod.spanID.Add(1)),
		}

		newCtx := ctx.WithParent(context.WithValue(ctx, logKey, ctxMod))

		ctxMod.log.Info("Receive")
		defer ctxMod.log.Info("Finish")

		newCtx.Next()

		return nil
	})

	return nil
}

// Register endpoints and launch the server
func LaunchWithLog() (addr string, cancel func()) {
	server, _ := espresso.New(withLogModule)

	server.HandleFunc(func(ctx espresso.Context) error {
		if err := ctx.Endpoint(http.MethodGet, "/").End(); err != nil {
			return err
		}

		logger := ctx.Value(logKey).(*logModule)
		logger.log.Info("In handle func.")
		return nil
	})

	httpSvr := httptest.NewServer(server)
	addr = httpSvr.URL
	cancel = func() {
		httpSvr.Close()
	}

	return
}

func ExampleCustomLog() {
	addr, cancel := LaunchWithLog()
	defer cancel()

	{
		resp, err := http.Get(addr + "/")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}
	{
		resp, err := http.Get(addr + "/")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}
	// Output:
	// {"level":"INFO","msg":"Receive","method":"GET","path":"/","span":1}
	// {"level":"INFO","msg":"In handle func.","method":"GET","path":"/","span":1}
	// {"level":"INFO","msg":"Finish","method":"GET","path":"/","span":1}
	// {"level":"INFO","msg":"Receive","method":"GET","path":"/","span":2}
	// {"level":"INFO","msg":"In handle func.","method":"GET","path":"/","span":2}
	// {"level":"INFO","msg":"Finish","method":"GET","path":"/","span":2}
}
