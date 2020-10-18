package metric

import "github.com/prometheus/client_golang/prometheus"

// vectorOpts contains the common arguments for creating vec Metric..
type vectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type CounterVecOpts vectorOpts

type CounterVec interface {
	Inc(labels ...string)
	Add(v float64, labels ...string)
}

type promCounterVec struct {
	counter *prometheus.CounterVec
}

func NewCounterVec(cfg *CounterVecOpts) CounterVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promCounterVec{
		counter: vec,
	}
}

func (vec *promCounterVec) Inc(labels ...string) {
	vec.counter.WithLabelValues(labels...).Inc()
}

func (vec *promCounterVec) Add(v float64, labels ...string) {
	vec.counter.WithLabelValues(labels...).Add(v)
}

