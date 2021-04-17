package mem

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/log"
	"time"
)


func NewBigCache(maxSize int) *bc {
	c, err := bigcache.NewBigCache(bigcache.Config{
		HardMaxCacheSize: maxSize,
	})
	if err != nil {
		log.Fatal("cache init failed: %s", err.Error())
	}
	log.Info("[mem bigcache] init success")
	return &bc{
		c,
	}
}

type bc struct {
	cache *bigcache.BigCache
}

func (c *bc) Get(key string, dest interface{}) error {
	dest, err := c.cache.Get(key)
	if err != nil && errors.Is(err, bigcache.ErrEntryNotFound) {
		return cache.ErrNotFound
	}
	return nil
}

func (c *bc) Set(k string, v interface{}) error {
	bs, ok := v.([]byte)
	if !ok {
		return errors.New("need bytes")
	}
	return c.cache.Set(k, bs)
}

func (c *bc) Del(keys ...string) error {
	for _, k:= range keys {
		if err := c.cache.Delete(k); err != nil {
			return fmt.Errorf("delete key: %s, failed: %s", k, err.Error())
		}
	}
	return nil
}

func (c *bc) Take(dest interface{}, key string, query func(v interface{}) error) error {
	err := c.Get(key, dest)
	if errors.Is(cache.ErrNotFound, err) {
		if err := query(dest); err != nil {
			return err
		}
		_ = c.Set(key, dest)
	}
	return nil
}

func (c *bc) SetWithExpire(key string, v interface{}, expire time.Duration) error {
	log.Warn("mem bigcache do not support set with expire")
	return c.Set(key, v)
}

func (c *bc) TakeWithExpire(dest interface{}, key string, query func(v interface{}) error, expire time.Duration) error {
	log.Warn("mem bigcache do not support take with expire")
	return c.Take(dest, key, query)
}

