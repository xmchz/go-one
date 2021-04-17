package storage

import (
	"context"
	"errors"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/log"
)

const CtxCacheKey = "storage-cache-key"

var ErrCacheKeyNotExist = errors.New("cache key not exist in context")

type Cache struct {
	Storage
	cache.Cache
}

func (s *Cache) key(ctx context.Context) (string, error) {
	key, ok := ctx.Value(CtxCacheKey).(string)
	if !ok {
		log.Debug("ctx cache key not exist")
		return "", ErrCacheKeyNotExist
	}
	return key, nil
}

func (s *Cache) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	key, err := s.key(ctx)
	if err != nil {
		return s.Storage.Find(ctx, dest, query, args...)
	}
	return s.Cache.Take(dest, key, func(v interface{}) error {
		return s.Storage.Find(ctx, dest, query, args...)
	})
}

func (s *Cache) Update(ctx context.Context, query string, args ...interface{}) error {
	key, err := s.key(ctx)
	if err != nil {
		return s.Storage.Update(ctx, query, args...)
	}
	if err := s.Storage.Update(ctx, query, args...); err != nil {
		return err
	}
	return s.Cache.Del(key)
}

func (s *Cache) Delete(ctx context.Context, query string, args ...interface{}) error {
	key, err := s.key(ctx)
	if err != nil {
		return s.Storage.Delete(ctx, query, args...)
	}
	if err := s.Storage.Delete(ctx, query, args...); err != nil {
		return err
	}
	return s.Cache.Del(key)
}
