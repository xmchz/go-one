// global Prometheus metrics registry.
package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xmchz/go-one/metric"
)

type Counter struct {
	cv  *prometheus.CounterVec
	lvs metric.LabelValues
}

func NewCounter(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return &Counter{
		cv: cv,
	}
}

func (m *Counter) With(labelValues ...string) metric.Counter {
	return &Counter{
		cv:  m.cv,
		lvs: m.lvs.With(labelValues...),
	}
}

func (m *Counter) Add(delta float64) {
	m.cv.With(toPromLabels(m.lvs...)).Add(delta)
}

func toPromLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
