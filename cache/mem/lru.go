package mem

import (
	"errors"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/container"
	"github.com/xmchz/go-one/log"
	"sync"
	"time"
)

func New() *lru {
	return &lru{
		Lru: container.NewLru(3),
	}
}

type lru struct {
	*container.Lru
	mu sync.Mutex
}

func (c *lru) Get(key string, dest interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var ok bool
	if dest, ok = c.Lru.Get(key); !ok {
		log.Debug("mem cache miss, key: %s", key)
		return cache.ErrNotFound
	}
	log.Debug("mem cache hit, key: %s", key)
	return nil
}

func (c *lru) Set(key string, val interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Lru.Set(key, val)
	return nil
}

func (c *lru) Del(keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		log.Debug("mem cache delete, key: %s", key)
		c.Lru.Del(key)
	}
	return nil
}

func (c *lru) Take(dest interface{}, key string, query func(v interface{}) error) error {
	err := c.Get(key, dest)
	if errors.Is(cache.ErrNotFound, err) {
		if err := query(dest); err != nil {
			return err
		}
		_ = c.Set(key, dest)
	}
	return nil
}

func (c *lru) SetWithExpire(key string, v interface{}, expire time.Duration) error {
	return c.Set(key, v)
}

func (c *lru) TakeWithExpire(dest interface{}, key string, query func(v interface{}) error, expire time.Duration) error {
	return c.Take(dest, key, query)
}
