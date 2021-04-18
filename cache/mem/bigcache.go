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
		log.Fatal("mem bigcache init failed: %s", err.Error())
	}
	log.Info("[mem bigcache] init success")
	return &bc{
		c,
	}
}

type bc struct {
	cache *bigcache.BigCache
}

func (c *bc) Get(k string, dest interface{}) error {
	bs, err := c.cache.Get(k)
	if err != nil && errors.Is(err, bigcache.ErrEntryNotFound) {
		log.Debug("mem bigcache miss, key: %s", k)
		return cache.ErrNotFound
	}
	err = json.Unmarshal(bs, &dest)
	if err != nil {
		return fmt.Errorf("mem bigcache unmarshal key:%s, err:%w", k, err)
	}
	log.Debug("mem bigcache hit, key: %s", k)
	return nil
}

func (c *bc) Set(k string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("mem bigcache marshal key:%s, err:%w", k, err)
	}
	return c.cache.Set(k, bs)
}

func (c *bc) Del(keys ...string) error {
	for _, k:= range keys {
		if err := c.cache.Delete(k); err != nil {
			return fmt.Errorf("mem bigcache delete key: %s, err: %w", k, err)
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
			log.Error("mem bigcache take key:%s, err:%s", k, err.Error())
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

