package storage

import (
	"context"
	"database/sql"
	"github.com/xmchz/go-one/breaker"
)

type breakerStorage struct {
	Storage
	breaker.Breaker
}

func (s *breakerStorage) acceptable(err error) bool {
	return err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone
}

func (s *breakerStorage) Create(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {

	return res, s.Do(func() error {
		res, err  = s.Storage.Create(ctx, query, args...)
		return err
	}, s.acceptable)
}

func (s *breakerStorage) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Find(ctx, dest, query, args...)
	}, s.acceptable)
}

func (s *breakerStorage) FindList(ctx context.Context, dests interface{}, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.FindList(ctx, dests, query, args...)
	}, s.acceptable)
}

func (s *breakerStorage) FindListIn(ctx context.Context, dests interface{}, query string, set interface{}) error {
	return s.Do(func() error {
		return s.Storage.FindListIn(ctx, dests, query, set)
	}, s.acceptable)
}

func (s *breakerStorage) Update(ctx context.Context, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Update(ctx, query, args...)
	}, s.acceptable)
}

func (s *breakerStorage) Delete(ctx context.Context, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Delete(ctx, query, args...)
	}, s.acceptable)
}
