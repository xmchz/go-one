package storage

import (
	"github.com/xmchz/go-one/breaker"
	"github.com/xmchz/go-one/cache"
)

type Option func(Storage) Storage

func WithBreaker(breaker breaker.Breaker) Option {
	return func(s Storage) Storage{
		return &Breaker{
			Storage: s,
			Breaker: breaker,
		}
	}
}

func WithCache(cache cache.Cache) Option {
	return func(s Storage) Storage{
		return &Cache{
			Storage: s,
			Cache:   cache,
		}
	}
}
