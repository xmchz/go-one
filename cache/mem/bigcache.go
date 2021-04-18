package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/log"
	"time"
)


func NewBigCache(maxSize int, eviction time.Duration) *bc {
	config := bigcache.DefaultConfig(eviction)
	config.HardMaxCacheSize = maxSize
	c, err := bigcache.NewBigCache(config)
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
		log.Debug("mem bigcache miss, key: %s", key)
		return cache.ErrNotFound
	}
	log.Debug("mem bigcache hit, key: %s", key)
	return nil
}

func (c *bc) Set(k string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("cache marshal key:%s, err:%w", k, err)
	}
	return c.cache.Set(k, bs)
}

func (c *bc) Del(keys ...string) error {
	for _, k:= range keys {
		if err := c.cache.Delete(k); err != nil {
			return fmt.Errorf("delete key: %s, err: %w", k, err)
		}
	}
	return nil
}

func (c *bc) Take(dest interface{}, k string, query func(v interface{}) error) error {
	err := c.Get(k, dest)
	if errors.Is(cache.ErrNotFound, err) {
		if err := query(dest); err != nil {
			return err
		}
		err = c.Set(k, dest)
		if err != nil {
			log.Error("cache take key:%s, err:%s", k, err.Error())
		}
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

