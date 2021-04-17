package cache

import (
	"errors"
	"github.com/xmchz/go-one/metric"
	"sync/atomic"
)

type Metric struct {
	Cache
	metric.Counter
}

func (c *Metric) Get(key string, dest interface{}) (err error) {
	err = c.Cache.Get(key, dest)
	if err != nil {
		c.With("cache", "miss").Add(1)
	} else {
		c.With("cache", "hit").Add(1)
	}
	return
}

func (c *Metric) Take(dest interface{}, key string, query func(v interface{}) error) error {
	err := c.Get(key, dest)
	if errors.Is(ErrNotFound, err) {
		if err := query(dest); err != nil {
			return err
		}
		_ = c.Set(key, dest)
	}
	return nil
}

type Stats struct {
	Hit   uint64
	Miss  uint64
	cur   string
}

func (m *Stats) With(labelValues ...string) metric.Counter {
	if len(labelValues) != 2 {
		return m
	}
	m.cur = labelValues[1]
	return m
}

func (m *Stats) Add(delta float64) {
	switch m.cur {
	case "miss":
		atomic.AddUint64(&m.Miss, uint64(delta))
	case "hit":
		atomic.AddUint64(&m.Hit, uint64(delta))
	}
}

func (m *Stats) Value() (uint64, uint64) {
	return atomic.LoadUint64(&m.Hit), atomic.LoadUint64(&m.Miss)
}
