package metric

import "github.com/prometheus/client_golang/prometheus"


type GaugeVecOpts vectorOpts

// GaugeVec gauge vec.
type GaugeVec interface {
	Set(v float64, labels ...string)
	Inc(labels ...string)
	Add(v float64, labels ...string)
}

type promGaugeVec struct {
	gauge *prometheus.GaugeVec
}

func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promGaugeVec{
		gauge: vec,
	}
}

// Inc Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (vec *promGaugeVec) Inc(labels ...string) {
	vec.gauge.WithLabelValues(labels...).Inc()
}

// Add Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (vec *promGaugeVec) Add(v float64, labels ...string) {
	vec.gauge.WithLabelValues(labels...).Add(v)
}

// Set set the given value to the collection.
func (vec *promGaugeVec) Set(v float64, labels ...string) {
	vec.gauge.WithLabelValues(labels...).Set(v)
}
