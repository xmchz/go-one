package cache

import "github.com/xmchz/go-one/metric"

type Option func(Cache) Cache

func WithBlock() Option {
	return func(cache Cache) Cache {
		return &Block{
			Cache: cache,
		}
	}
}

func WithMetric(counter metric.Counter) Option {
	return func(cache Cache) Cache {
		return &Metric{
			Cache: cache,
			Counter: counter,
		}
	}
}
