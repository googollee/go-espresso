package main

import (
	"context"
	"net/http"
	"os"

	"github.com/googollee/go-espresso"
	"github.com/googollee/go-espresso/logger"
	"github.com/googollee/go-espresso/prometheus"
	"google.golang.org/api/tracing/v2"
)

type Data struct{}

type Service struct {
	logger logger.Logger
	metric metric.Int
	tracer tracing.Tracer
}

func (s *Service) HandleCreate(ctx espresso.Context[Data]) error {
	if err := ctx.Endpoint(http.MethodPost, "/service").End(); err != nil {
		return espresso.ErrWithStatus(http.StatusBadRequest, err)
	}

	return espresso.Procedure(ctx, s.Create)
}

func (s *Service) Create(ctx context.Context, arg int) (string, error) {
	s.metric.Add(1)

	ctx, span := s.tracer.Start(ctx)
	defer span.End()

	s.logger.Info(ctx, "in create")
	s.logger.Debug(ctx, "input", "arg", arg)

	req := http.NewRequest()
	s.tracer.Inject(ctx, req.Header)
	resp := http.DefaultClient.Do(req)

	return "", nil
}

func main() {
	svc := &Service{}

	eng := espresso.NewEngine(
		logger.New(os.Stderr, logger.DEBUG),
		prometheus.New("/metrics"),
		tracing.New("https://url"),
	)

	svc.logger = eng.Logger()
	svc.metric = prometheus.IntGauge(eng, "calls")
	svc.tracer = eng.Tracer()

	eng.HandleAll(svc)

	eng.ListenAndServe(":8080")
}
