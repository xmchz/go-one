package mem

import (
	"github.com/xmchz/go-one/metric"
	"sync/atomic"
)

type Counter struct {
	total uint64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (m *Counter) With(labelValues ...string) metric.Counter {
	return m
}

func (m *Counter) Add(delta float64) {
	atomic.AddUint64(&m.total, uint64(delta))
}

func (m *Counter) Value() uint64 {
	return atomic.LoadUint64(&m.total)
}
