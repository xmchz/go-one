package storage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/xmchz/go-one/log"
)

const (
	ctxTxKey = "context-tx-key"
)

func New(conf Config, opts ...Option) Storage {
	pool, err := sqlx.Connect(conf.DbDriver(), conf.DbUrl())
	if err != nil {
		log.Fatal("connect DB failed, err:%v", err)
	}
	pool.SetMaxOpenConns(conf.MaxOpenConn())
	pool.SetMaxIdleConns(conf.MaxIdleConn())
	var s Storage
	s = &storage{pool}
	for _, opt := range opts {
		s = opt(s)
	}
	return s
}

type Storage interface {
	BeginWithTx(ctx context.Context) (context.Context, func(err error), error)
	Create(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Find(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	FindList(ctx context.Context, dests interface{}, query string, args ...interface{}) error
	FindListIn(ctx context.Context, dests interface{}, query string, set interface{}) error
	Update(ctx context.Context, query string, args ...interface{}) error
	Delete(ctx context.Context, query string, args ...interface{}) error
}

type storage struct {
	*sqlx.DB
}

func (s *storage) BeginWithTx(ctx context.Context) (context.Context, func(err error), error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	return context.WithValue(ctx, ctxTxKey, tx), func(err error) {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
		if err != nil {
			log.Debug("rollback")
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			log.Debug("commit")
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}, err
}

func (s *storage) Create(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.exec(ctx, query, args...)
}

func (s *storage) Update(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.exec(ctx, query, args...)
	return err
}

func (s *storage) Delete(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.exec(ctx, query, args...)
	return err
}

func (s *storage) Find(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	log.Debug("query from db")
	if tx := s.tx(ctx); tx != nil {
		err = tx.GetContext(ctx, dest, query, args...)
		return
	}
	err = s.DB.GetContext(ctx, dest, query, args...)
	return
}

func (s *storage) FindList(ctx context.Context, dests interface{}, query string, args ...interface{}) (err error) {
	if tx := s.tx(ctx); tx != nil {
		err = tx.SelectContext(ctx, dests, query, args...)
		return
	}
	err = s.DB.SelectContext(ctx, dests, query, args...)
	return
}

func (s *storage) FindListIn(ctx context.Context, dests interface{}, query string, set interface{}) (err error) {
	var args []interface{}
	query, args, err = sqlx.In(query, set)
	if err != nil {
		return
	}
	if tx := s.tx(ctx); tx != nil {
		err = tx.SelectContext(ctx, dests, query, args...)
		return
	}
	err = s.DB.SelectContext(ctx, dests, s.DB.Rebind(query), args...)
	return
}

func (s *storage) exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := s.tx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return s.DB.ExecContext(ctx, query, args...)
}

func (s *storage) tx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(ctxTxKey).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}
