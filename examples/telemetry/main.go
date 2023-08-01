//go:build wip

package main

import (
	"net/http"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/log"
	prometheus "github.com/googollee/go-espresso/monitoring/prometheus"
	openapi "github.com/googollee/go-espresso/openapi"
	tracing "github.com/googollee/go-espresso/tracing"
	opentelemetry "github.com/googollee/go-espresso/tracing/opentelemetry"
)

var rpcCalls = prometheus.IntGauge("rpcCalls")

type Data struct{}

type Service struct{}

func (s *Service) Create(ctx espresso.Context) error {
	return espresso.Produce(ctx, s.create)
}

func (s *Service) create(ctx espresso.Context, arg int) (string, error) {
	if err := ctx.Endpoint(http.MethodPost, "/service").End(); err != nil {
		return "", espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	rpcCalls.Add(1)

	ctx, done := tracing.Start(ctx)
	defer done()

	log.Info(ctx, "in create")
	log.Debug(ctx, "input", "arg", arg)

	// or:
	// req, err := tracing.NewHTTPRequest(http.MethodPost, "", nil)
	req, _ := http.NewRequest(http.MethodPost, "", nil)
	tracing.Inject(ctx, req.Header)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	return "", nil
}

func main() {
	svc := &Service{}

	svr, _ := espresso.New()
	svr.With(
		prometheus.New("/metrics"),
		opentelemetry.New("https://url"),
		openapi.New("/spec"),
	)

	svr.HandleAll(svc)
	svr.ListenAndServe(":8080")
}
