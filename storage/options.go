package storage

import (
	"github.com/xmchz/go-one/breaker"
	"github.com/xmchz/go-one/cache"
)

type Option func(Storage) Storage

func WithBreaker(breaker breaker.Breaker) Option {
	return func(s Storage) Storage {
		return &breakerStorage{
			Storage: s,
			Breaker: breaker,
		}
	}
}

func WithCache(cache cache.Cache) Option {
	return func(s Storage) Storage {
		return &cacheStorage{
			Storage: s,
			Cache:   cache,
		}
	}
}
