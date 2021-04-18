package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/log"
)

const (
	ctxKeyCache       = "storage-cache-set-key"    // set cache
	ctxKeyDeleteCache = "storage-cache-delete-key" // delete cache
)

var ErrCacheKeyNotExist = errors.New("cache key not exist in context")

func CtxSetKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxKeyCache, key)
}

func CtxDeleteKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxKeyDeleteCache, key)
}

type Cache struct {
	Storage
	cache.Cache
}

func (s *Cache) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	key, err := s.key(ctx, ctxKeyCache)
	if errors.Is(err, ErrCacheKeyNotExist) {
		return s.Storage.Find(ctx, dest, query, args...)
	}
	return s.Cache.Take(dest, key, func(v interface{}) error {
		log.Debug("query from db, key:%s", key)
		return s.Storage.Find(ctx, dest, query, args...)
	})
}

func (s *Cache) FindList(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	key, err := s.key(ctx, ctxKeyCache)
	if errors.Is(err, ErrCacheKeyNotExist) {
		return s.Storage.FindList(ctx, dest, query, args...)
	}
	return s.Cache.Take(dest, key, func(v interface{}) error {
		log.Debug("query from db, key:%s", key)
		return s.Storage.FindList(ctx, dest, query, args...)
	})
}

func (s *Cache) Create(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	res, err := s.Storage.Create(ctx, query, args...)
	if err != nil {
		return res, err
	}
	s.deleteCtxCache(ctx)
	return res, nil
}

func (s *Cache) Update(ctx context.Context, query string, args ...interface{}) error {
	if err := s.Storage.Update(ctx, query, args...); err != nil {
		return err
	}
	if key, err := s.key(ctx, ctxKeyCache); err == nil {
		s.deleteCache(key)
	}
	s.deleteCtxCache(ctx)
	return nil
}

func (s *Cache) Delete(ctx context.Context, query string, args ...interface{}) error {
	if err := s.Storage.Delete(ctx, query, args...); err != nil {
		return err
	}
	if key, err := s.key(ctx, ctxKeyCache); err == nil {
		s.deleteCache(key)
	}
	s.deleteCtxCache(ctx)
	return nil
}

func (s *Cache) key(ctx context.Context, ctxKey string) (string, error) {
	key, ok := ctx.Value(ctxKey).(string)
	if !ok {
		// log.Debug("%s not exist in ctx", ctxKey)
		return "", ErrCacheKeyNotExist
	}
	return key, nil
}

func (s *Cache) deleteCache(k string) {
	err := s.Cache.Del(k)
	if err != nil {
		log.Error("storage delete cache key:%s, err:%s", k, err.Error())
	}
}

func (s *Cache) deleteCtxCache(ctx context.Context) {
	if k, err := s.key(ctx, ctxKeyDeleteCache); err == nil {
		s.deleteCache(k)
	}
}
