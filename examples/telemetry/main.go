//go:build wip

package main

import (
	"context"
	"net/http"

	"github.com/googollee/go-espresso"
	log "github.com/googollee/go-espresso/log"
	prometheus "github.com/googollee/go-espresso/monitoring/prometheus"
	openapi "github.com/googollee/go-espresso/openapi"
	tracing "github.com/googollee/go-espresso/tracing"
	opentelemetry "github.com/googollee/go-espresso/tracing/opentelemetry"
	"golang.org/x/exp/slog"
)

var rpcCalls = prometheus.IntGauge("rpcCalls")

type Data struct{}

type Service struct{}

func (s *Service) Create(ctx espresso.Context[Data]) error {
	return espresso.Procedure(ctx, s.Create)
}

func (s *Service) create(ctx context.Context, arg int) (string, error) {
	if err := ctx.Endpoint(http.MethodPost, "/service").End(); err != nil {
		return "", espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	rpcCalls.Add(1)

	ctx, span := tracing.Start(ctx)
	defer span.End()

	log.Info(ctx, "in create")
	log.Debug(ctx, "input", "arg", arg)

	// or:
	// req, err := tracing.NewHTTPRequest(http.MethodPost, "", nil)
	req, err := http.NewRequest(http.MethodPost, "", nil)
	tracing.Inject(ctx, req.Header)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	return "", nil
}

func main() {
	svc := &Service{}

	svr := espresso.Default().With(
		log.New(slog.Default()),
		prometheus.New("/metrics"),
		opentelemetry.New("https://url"),
		openapi.New("/spec"),
	)

	svr.HandleAll(svc)
	svr.ListenAndServe(":8080")
}
