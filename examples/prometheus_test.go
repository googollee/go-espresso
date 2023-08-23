package espresso_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/monitoring/prometheus"
)

var mycounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "mycounter",
})

func AddCounter(ctx espresso.Context) error {
	var num int
	if err := ctx.Endpoint(http.MethodGet, "/inc/:num").
		BindPath("num", &num).End(); err != nil {
		return err
	}

	mycounter.Add(float64(num))

	ctx.ResponseWriter().WriteHeader(http.StatusNoContent)
	return nil
}

func LaunchWithPrometheus() (addr string, cancel func()) {
	server, _ := espresso.New(prometheus.New("/metrics"))

	server.HandleFunc(AddCounter)

	httpSvr := httptest.NewServer(server)
	addr = httpSvr.URL
	cancel = func() {
		httpSvr.Close()
	}

	return
}

func ExampleMonitoringWithPrometheus() {
	addr, cancel := LaunchWithPrometheus()
	defer cancel()

	{
		resp, err := http.Get(addr + "/metrics")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
	}

	{
		resp, err := http.Get(addr + "/inc/100")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	{
		resp, err := http.Get(addr + "/metrics")
		if err != nil {
			panic(err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
	}

	// Output:
	// 200 text/plain; version=0.0.4; charset=utf-8 # HELP mycounter
	// # TYPE mycounter counter
	// mycounter 0
	// # HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
	// # TYPE promhttp_metric_handler_errors_total counter
	// promhttp_metric_handler_errors_total{cause="encoding"} 0
	// promhttp_metric_handler_errors_total{cause="gathering"} 0

	// 200 text/plain; version=0.0.4; charset=utf-8 # HELP mycounter
	// # TYPE mycounter counter
	// mycounter 100
	// # HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
	// # TYPE promhttp_metric_handler_errors_total counter
	// promhttp_metric_handler_errors_total{cause="encoding"} 0
	// promhttp_metric_handler_errors_total{cause="gathering"} 0
}
