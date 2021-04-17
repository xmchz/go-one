package storage

import (
	"context"
	"database/sql"
	"github.com/xmchz/go-one/breaker"
)

type Breaker struct {
	Storage
	breaker.Breaker
}

func (s *Breaker) acceptable(err error) bool {
	return err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone
}

func (s *Breaker) Create(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {

	return res, s.Do(func() error {
		res, err  = s.Storage.Create(ctx, query, args...)
		return err
	}, s.acceptable)
}

func (s *Breaker) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Find(ctx, dest, query, args...)
	}, s.acceptable)
}

func (s *Breaker) FindList(ctx context.Context, dests interface{}, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.FindList(ctx, dests, query, args...)
	}, s.acceptable)
}

func (s *Breaker) FindListIn(ctx context.Context, dests interface{}, query string, set interface{}) error {
	return s.Do(func() error {
		return s.Storage.FindListIn(ctx, dests, query, set)
	}, s.acceptable)
}

func (s *Breaker) Update(ctx context.Context, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Update(ctx, query, args...)
	}, s.acceptable)
}

func (s *Breaker) Delete(ctx context.Context, query string, args ...interface{}) error {
	return s.Do(func() error {
		return s.Storage.Delete(ctx, query, args...)
	}, s.acceptable)
}
