package espresso_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/googollee/go-espresso"
)

// Use `fmt.Print` to output logs
// `go test` replaces `fmt.Print` to verify the output.
type Printer struct{}

func (p Printer) Write(b []byte) (int, error) {
	return fmt.Print(string(b))
}

func newFmtLogger() *slog.Logger {
	w := Printer{}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove `time` attribute.
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(handler)
}

// Define a custom ResponseWriter to get the response status code.
type ResponseWriter struct {
	http.ResponseWriter
	Code int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Code = code
}

// Define a middleware to log requests.
func LogRequest(ctx espresso.Context) error {
	// Update ResponseWriter to record status code
	respW := &ResponseWriter{
		ResponseWriter: ctx.ResponseWriter(),
	}
	newCtx := espresso.WithResponseWriter(ctx, respW)

	// Set attributes for this request
	req := ctx.Request()
	newCtx = espresso.WithLogAttr(newCtx, "method", req.Method, "path", req.URL.Path)

	// Log the request
	espresso.Info(newCtx, "Request")

	// Log the response status code after all handlers
	defer func() {
		espresso.Info(newCtx, "Response", "code", respW.Code)
	}()

	// Call following handlers
	newCtx.Next()

	return nil
}

func LaunchLogServer() (addr string, cancel func()) {
	// Create a server with out default middlewares, and with a custom logger
	svr, _ := espresso.New(espresso.WithoutDefault(), espresso.WithLog(newFmtLogger()))

	// Add `LogRequest` as the first middlware
	svr.Use(LogRequest)

	// Define an endpoint
	svr.HandleFunc(func(ctx espresso.Context) error {
		if err := ctx.Endpoint(http.MethodGet, "/").End(); err != nil {
			return err
		}
		espresso.Info(ctx, "in function")
		ctx.ResponseWriter().WriteHeader(http.StatusNoContent)
		return nil
	})

	// Launch server
	tsvr := httptest.NewServer(svr)
	return tsvr.URL, func() { tsvr.Close() }
}

func ExampleLogRequest() {
	addr, cancel := LaunchLogServer()
	defer cancel()

	// The log records this request and its response.
	resp, _ := http.Get(addr)
	resp.Body.Close()

	// Output:
	// level=INFO msg=Register method=GET path=/ handler=espresso.HandleFunc
	// level=INFO msg=Request method=GET path=/
	// level=INFO msg="in function" method=GET path=/
	// level=INFO msg=Response method=GET path=/ code=204
}
