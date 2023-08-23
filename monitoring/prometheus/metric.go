package prometheus

import (
	"net/http"

	"github.com/googollee/go-espresso"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var DefaultRegistry = prometheus.NewRegistry()

type GaugeOpts = prometheus.GaugeOpts

func NewGauge(opt GaugeOpts) prometheus.Gauge {
	ret := prometheus.NewGauge(opt)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewGaugeFunc(opt GaugeOpts, fn func() float64) prometheus.GaugeFunc {
	ret := prometheus.NewGaugeFunc(opt, fn)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewGaugeVec(opt GaugeOpts, labels []string) *prometheus.GaugeVec {
	ret := prometheus.NewGaugeVec(opt, labels)
	DefaultRegistry.MustRegister(ret)
	return ret
}

type CounterOpts = prometheus.CounterOpts

func NewCounter(opt CounterOpts) prometheus.Counter {
	ret := prometheus.NewCounter(opt)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewCounterFunc(opt CounterOpts, fn func() float64) prometheus.CounterFunc {
	ret := prometheus.NewCounterFunc(opt, fn)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewCounterVec(opt CounterOpts, labels []string) *prometheus.CounterVec {
	ret := prometheus.NewCounterVec(opt, labels)
	DefaultRegistry.MustRegister(ret)
	return ret
}

type SummaryOpts = prometheus.SummaryOpts

func NewSummary(opt SummaryOpts) prometheus.Summary {
	ret := prometheus.NewSummary(opt)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewSummaryVec(opt SummaryOpts, labels []string) *prometheus.SummaryVec {
	ret := prometheus.NewSummaryVec(opt, labels)
	DefaultRegistry.MustRegister(ret)
	return ret
}

type HistogramOpts = prometheus.HistogramOpts

func NewHistogram(opt HistogramOpts) prometheus.Histogram {
	ret := prometheus.NewHistogram(opt)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func NewHistogramVec(opt HistogramOpts, labels []string) *prometheus.HistogramVec {
	ret := prometheus.NewHistogramVec(opt, labels)
	DefaultRegistry.MustRegister(ret)
	return ret
}

type UntypedOpts = prometheus.UntypedOpts

func NewUntypedFunc(opt UntypedOpts, fn func() float64) prometheus.UntypedFunc {
	ret := prometheus.NewUntypedFunc(opt, fn)
	DefaultRegistry.MustRegister(ret)
	return ret
}

func New(path string) espresso.ServerOption {
	handler := promhttp.HandlerFor(DefaultRegistry, promhttp.HandlerOpts{
		Registry: DefaultRegistry,
	})

	return func(s *espresso.Server) error {
		s.HandleFunc(func(ctx espresso.Context) error {
			if err := ctx.Endpoint(http.MethodGet, path).End(); err != nil {
				return err
			}

			handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
			return nil
		})

		return nil
	}
}
