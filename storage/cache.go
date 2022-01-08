package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/log"
)

type ctxKey uint

const (
	ctxKeyForCache ctxKey = 0
)

var ErrCacheKeyNotExist = errors.New("cache key not exist in context")

func CtxWithCacheKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxKeyForCache, key)
}

type cacheStorage struct {
	Storage
	cache.Cache
}

func (s *cacheStorage) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	cacheKey, ok := ctx.Value(ctxKeyForCache).(string)
	if !ok {
		return s.Storage.Find(ctx, dest, query, args...)
	}
	return s.Cache.Take(dest, cacheKey, func(_ interface{}) error {
		log.Debug("query from db, key:%s", cacheKey)
		return s.Storage.Find(ctx, dest, query, args...)
	})
}

func (s *cacheStorage) FindList(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	cacheKey, ok := ctx.Value(ctxKeyForCache).(string)
	if !ok {
		return s.Storage.FindList(ctx, dest, query, args...)
	}
	return s.Cache.Take(dest, cacheKey, func(_ interface{}) error {
		log.Debug("query from db, key:%s", cacheKey)
		return s.Storage.FindList(ctx, dest, query, args...)
	})
}

func (s *cacheStorage) Create(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	res, err := s.Storage.Create(ctx, query, args...)
	if err != nil {
		return res, err
	}
	s.deleteCache(ctx)
	return res, nil
}

func (s *cacheStorage) Update(ctx context.Context, query string, args ...interface{}) error {
	if err := s.Storage.Update(ctx, query, args...); err != nil {
		return err
	}
	s.deleteCache(ctx)
	return nil
}

func (s *cacheStorage) Delete(ctx context.Context, query string, args ...interface{}) error {
	if err := s.Storage.Delete(ctx, query, args...); err != nil {
		return err
	}
	s.deleteCache(ctx)
	return nil
}

func (s *cacheStorage) deleteCache(ctx context.Context) {
	cacheKey, ok := ctx.Value(ctxKeyForCache).(string)
	if !ok {
		return
	}
	err := s.Cache.Del(cacheKey)
	if err != nil {
		log.Error("storage delete cache key:%s, err:%s", cacheKey, err.Error())
	}
}
