package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type GaugeOpts = prometheus.GaugeOpts

func NewGauge(opt GaugeOpts) prometheus.Gauge {
	ret := prometheus.NewGauge(opt)
	prometheus.MustRegister(ret)
	return ret
}

func NewGaugeFunc(opt GaugeOpts, fn func() float64) prometheus.GaugeFunc {
	ret := prometheus.NewGaugeFunc(opt, fn)
	prometheus.MustRegister(ret)
	return ret
}

func NewGaugeVec(opt GaugeOpts, labels []string) *prometheus.GaugeVec {
	ret := prometheus.NewGaugeVec(opt, labels)
	prometheus.MustRegister(ret)
	return ret
}

type CounterOpts = prometheus.CounterOpts

func NewCounter(opt CounterOpts) prometheus.Counter {
	ret := prometheus.NewCounter(opt)
	prometheus.MustRegister(ret)
	return ret
}

func NewCounterFunc(opt CounterOpts, fn func() float64) prometheus.CounterFunc {
	ret := prometheus.NewCounterFunc(opt, fn)
	prometheus.MustRegister(ret)
	return ret
}

func NewCounterVec(opt CounterOpts, labels []string) *prometheus.CounterVec {
	ret := prometheus.NewCounterVec(opt, labels)
	prometheus.MustRegister(ret)
	return ret
}

type SummaryOpts = prometheus.SummaryOpts

func NewSummary(opt SummaryOpts) prometheus.Summary {
	ret := prometheus.NewSummary(opt)
	prometheus.MustRegister(ret)
	return ret
}

func NewSummaryVec(opt SummaryOpts, labels []string) *prometheus.SummaryVec {
	ret := prometheus.NewSummaryVec(opt, labels)
	prometheus.MustRegister(ret)
	return ret
}

type HistogramOpts = prometheus.HistogramOpts

func NewHistogram(opt HistogramOpts) prometheus.Histogram {
	ret := prometheus.NewHistogram(opt)
	prometheus.MustRegister(ret)
	return ret
}

func NewHistogramVec(opt HistogramOpts, labels []string) *prometheus.HistogramVec {
	ret := prometheus.NewHistogramVec(opt, labels)
	prometheus.MustRegister(ret)
	return ret
}

type UntypedOpts = prometheus.UntypedOpts

func NewUntypedFunc(opt UntypedOpts, fn func() float64) prometheus.UntypedFunc {
	ret := prometheus.NewUntypedFunc(opt, fn)
	prometheus.MustRegister(ret)
	return ret
}
