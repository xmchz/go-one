package cache

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("note found")
)

func New(cache Cache, opts ...Option) Cache {
	for _, opt := range opts {
		cache = opt(cache)
	}
	return cache
}

type Cache interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}) error
	Del(keys ...string) error
	Take(dest interface{}, key string, query func(v interface{}) error) error
	SetWithExpire(key string, v interface{}, expire time.Duration) error
	TakeWithExpire(dest interface{}, key string, query func(v interface{}) error, expire time.Duration) error
}
