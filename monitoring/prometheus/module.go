package prometheus

import (
	"context"
	"net/http"

	"github.com/googollee/go-espresso/basetype"
	"github.com/googollee/go-espresso/module"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Module = module.New[*Prometheus]()

func Use(opts ...Option) basetype.ServerOption {
	p := build(opts...)

	return func(s basetype.Server) error {
		s.HandleFunc(p.endpointMetric)
		return nil
	}
}

type Option func(*Prometheus)

func WithPath(path string) Option {
	return func(p *Prometheus) {
		if path == "" {
			return
		}

		p.path = path
	}
}

func WithRegistry(registry *prometheus.Registry) Option {
	return func(p *Prometheus) {
		if registry == nil {
			return
		}

		p.registry = registry
	}
}

func WithGatherers(gatherer ...prometheus.Gatherer) Option {
	return func(p *Prometheus) {
		p.gatherers = gatherer
	}
}

type Prometheus struct {
	path          string
	registry      prometheus.Registerer
	gatherers     prometheus.Gatherers
	metricHandler http.Handler
}

func build(opts ...Option) *Prometheus {
	ret := &Prometheus{
		path:      "/metrics",
		registry:  prometheus.DefaultRegisterer,
		gatherers: []prometheus.Gatherer{prometheus.DefaultGatherer},
	}

	for _, opt := range opts {
		opt(ret)
	}

	ret.metricHandler = promhttp.InstrumentMetricHandler(ret.registry, promhttp.HandlerFor(ret.gatherers, promhttp.HandlerOpts{}))

	return ret
}

func (p *Prometheus) CheckHealth(ctx context.Context) error {
	return nil
}

func (p *Prometheus) endpointMetric(ctx basetype.Context) error {
	if err := ctx.Endpoint(http.MethodGet, p.path).End(); err != nil {
		return err
	}

	p.metricHandler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	return nil
}
