package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmchz/go-one/breaker"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/cache/mem"
	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
	"github.com/xmchz/go-one/storage"
	"os"
	"sync"
	"testing"
)

const (
	sqlInsert     = `INSERT INTO wiki_content(content) VALUES(?)`
	sqlInsertErr  = `INSERT INTO wiki_content(id,content) VALUES(5,"")`
	sqlFind       = "SELECT content FROM wiki_content WHERE id=?"
	sqlUpdate     = "UPDATE wiki_content SET content=? WHERE id=?"
	sqlFindList   = "SELECT content FROM wiki_content WHERE id<?"
	sqlFindListIn = "SELECT content FROM wiki_content WHERE id in (?)"
)

type dbConf struct{}

func (d dbConf) DbDriver() string {
	return "mysql"
}

func (d dbConf) DbName() string {
	return "chz_wiki"
}

func (d dbConf) DbUrl() string {
	return "root:chouchou@tcp(127.0.0.1:3306)/chz_wiki?charset=utf8mb4&parseTime=True"
}

func (d dbConf) DbInfo() string {
	return "root@127.0.0.1:3306/chz_wiki"
}

func (d dbConf) MaxIdleConn() int {
	return 5
}

func (d dbConf) MaxOpenConn() int {
	return 5
}

func TestMain(m *testing.M) {
	log.Init(
		log.WithWriters(&writer.Console{Formatter: &formatter.Text{}}),
		log.WithLevel(log.DebugLevel),
	)
	exitVal := m.Run()
	log.Stop()
	os.Exit(exitVal)
}

func TestStorage_Create(t *testing.T) {
	s := storage.New(&dbConf{})
	_, err := s.Create(context.Background(), sqlInsert, "test content test content test content")
	assert.Nil(t, err)
}

func TestStorage_Tx(t *testing.T) {
	s := storage.New(&dbConf{})
	ctx, cleanTx, err := s.BeginWithTx(context.Background())
	if err != nil {
		return
	}
	defer func() { cleanTx(err) }()
	if _, err = s.Create(ctx, sqlInsert, "test content test content test content"); err != nil {
		return
	}
	if _, err = s.Create(ctx, sqlInsertErr); err != nil {
		return
	}
}

func TestStorage_Tx2(t *testing.T) {
	s := storage.New(&dbConf{})
	ctx, cleanTx, err := s.BeginWithTx(context.Background())
	if err != nil {
		return
	}
	defer func() { cleanTx(err) }()
	if _, err = s.Create(ctx, sqlInsert, "test content test content test content"); err != nil {
		return
	}
	if _, err = s.Create(ctx, sqlInsert, "test content test content test content"); err != nil {
		return
	}
}

func TestNew(t *testing.T) {
	s := storage.New(&dbConf{})
	var res string
	err := s.Find(context.Background(), &res, sqlFind, 5)
	require.Nil(t, err)
	require.NotEmpty(t, res)
}

func TestWithBreaker(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithBreaker(&breaker.Empty{}),
	)
	var res string
	err := s.Find(context.Background(), &res, sqlFind, 5)
	require.Nil(t, err)
	require.NotEmpty(t, res)
}

func TestWithCache(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithCache(mem.New()),
	)
	var res string
	for _, id := range []int{12, 13, 14, 12, 15, 13} {
		err := s.Find(
			context.WithValue(context.Background(), storage.CtxCacheKey, fmt.Sprintf("content-id-%d", id)),
			&res,
			sqlFind,
			id,
		)
		require.Nil(t, err)
	}
}

func TestWithCache2(t *testing.T) {
	s := storage.New(
		&dbConf{},
		//storage.WithCache(mem.New()),
		storage.WithCache(&cache.Block{
			Cache: mem.New(),
		}),
	)
	id := 12
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			var res string
			_ = s.Find(context.WithValue(context.Background(), storage.CtxCacheKey, fmt.Sprintf("content-id-%d", id)),
				&res,
				sqlFind,
				id,
			)
			wg.Done()
		}()
	}
	wg.Wait()
	var res string
	_ = s.Find(context.WithValue(context.Background(), storage.CtxCacheKey, fmt.Sprintf("content-id-%d", id)),
		&res,
		sqlFind,
		id,
	)
}

func TestStorage_Find(t *testing.T) {
	s := storage.New(
		&dbConf{},
	)
	var res string
	err := s.Find(context.Background(), &res, sqlFind, -1)
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestStorage_Update(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithCache(mem.New()),
	)
	var res string
	id := 12
	ctx := context.WithValue(context.Background(), storage.CtxCacheKey, fmt.Sprintf("content-id-%d", id))
	for i := 0; i < 2; i++ {
		err := s.Find(ctx, &res, sqlFind, id)
		require.Nil(t, err)
	}
	err := s.Update(ctx, sqlUpdate, "1"+res, id)
	require.Nil(t, err)
	err = s.Find(ctx, &res, sqlFind, id)
	require.Nil(t, err)
}

func TestStorage_FindList(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithCache(mem.New()),
	)
	var res []string
	err := s.FindList(context.Background(), &res, sqlFindList, 12)
	assert.Nil(t, err)
	assert.Len(t, res, 2)
}

func TestStorage_FindListIn(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithCache(mem.New()),
	)
	var res []string
	err := s.FindListIn(context.Background(), &res, sqlFindListIn, []int{12, 13})
	assert.Nil(t, err)
	assert.Len(t, res, 2)
}
