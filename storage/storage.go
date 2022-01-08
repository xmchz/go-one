package storage

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	sqlxAdapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/jmoiron/sqlx"
	"github.com/xmchz/go-one/access"
	"github.com/xmchz/go-one/log"
)

type txKey uint

const (
	ctxKeyForTx txKey = 0
)

func New(conf Config, opts ...Option) Storage {
	pool, err := sqlx.Connect(conf.DbDriver(), conf.DbUrl())
	if err != nil {
		log.Fatal("connect DB failed, err:%v", err)
	}
	log.Info("connect DB success: %s", conf.DbInfo())
	pool.SetMaxOpenConns(conf.MaxOpenConn())
	pool.SetMaxIdleConns(conf.MaxIdleConn())
	doMigrate(conf, pool.DB)
	var s Storage
	s = &storage{pool}
	for _, opt := range opts {
		s = opt(s)
	}
	log.Info("storage init success")
	return s
}

func doMigrate(conf Config, db *sql.DB) {
	if conf.MigrationUrl() == "" {
		log.Info("migrate skipped")
		return
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("migrate connect DB failed, err:%v", err)
	}
	defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance(
		conf.MigrationUrl(),
		conf.DbName(),
		driver,
	)
	if err != nil {
		log.Fatal("migrate err:%v", err)
	}
	err = m.Migrate(conf.MigrationVersion())
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate err:%v", err)
	} else if err == migrate.ErrNoChange {
		log.Info("migrate DB success: %d", conf.MigrationVersion())
	}
}

type Storage interface {
	access.Storage
	BeginWithTx(ctx context.Context) (context.Context, func(interface{}, error) error, error)
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

func (s *storage) AccessAdapter(tblName string) (access.Adapter, error) {
	return sqlxAdapter.NewAdapter(s.DB, tblName)
}

func (s *storage) BeginWithTx(ctx context.Context) (context.Context, func(recoveredPanic interface{}, err error) error, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func(recoveredPanic interface{} , err error) error {
		if recoveredPanic != nil {
			log.Debug("rollback by panic: %v", recoveredPanic)
			if err := tx.Rollback(); err != nil {
				log.Error("rollback err: %s", err)
				return err
			}
			log.Debug("rollback success")
			return nil
		}
		if err != nil {
			log.Debug("rollback by err: %s", err)
			if err := tx.Rollback(); err != nil {
				log.Error("rollback err: %s", err)
				return err
			}
			log.Debug("rollback success")
			return nil
		}
		if err := tx.Commit(); err != nil {
			log.Error("commit err: %s", err)
			return err
		}
		log.Debug("commit success")
		return nil
	}
	return context.WithValue(ctx, ctxKeyForTx, tx), cleanFunc, nil
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
	tx, ok := ctx.Value(ctxKeyForTx).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}
