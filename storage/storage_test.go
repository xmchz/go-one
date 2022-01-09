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
	logCore "github.com/xmchz/go-one/log/core"
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

func (d dbConf) MigrationUrl() string {
	return ""
}
func (d dbConf) MigrationVersion() uint {
	return 0
}

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
		logCore.WithWriters(writer.NewConsole(&formatter.Text{})),
		logCore.WithLevel(logCore.DebugLevel),
	)
	exitVal := m.Run()
	log.Stop()
	os.Exit(exitVal)
}

func TestStorage_Create(t *testing.T) {
	s := storage.New(&dbConf{})
	_, err := s.Create(context.Background(), sqlInsert, "test content23")
	assert.Nil(t, err)
}

func TestStorage_Tx(t *testing.T) {
	s := storage.New(&dbConf{})
	ctx, cleanTx, err := s.BeginWithTx(context.Background())
	require.Nil(t, err)
	defer func() {
		cleanErr := cleanTx(err, recover())
		require.Nil(t, cleanErr)
	}()
	if _, err = s.Create(ctx, sqlInsert, ""); err != nil {
		return
	}
	if _, err = s.Create(ctx, sqlInsertErr); err != nil {
		return
	}
}

func TestStorage_Tx2(t *testing.T) {
	s := storage.New(&dbConf{})
	ctx, cleanTx, err := s.BeginWithTx(context.Background())
	require.Nil(t, err)
	defer func() {
		cleanErr := cleanTx(err, recover())
		require.Nil(t, cleanErr)
	}()
	if _, err = s.Create(ctx, sqlInsert, ""); err != nil {
		return
	}
	if _, err = s.Create(ctx, sqlInsert, ""); err != nil {
		return
	}
	panic("panic in tx")
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
			storage.CtxWithCacheKey(context.Background(), fmt.Sprintf("content-id-%d", id)),
			&res,
			sqlFind,
			id,
		)
		require.Nil(t, err)
	}
}

func TestWithCacheConcurrency(t *testing.T) {
	c := cache.New(
		mem.New(),
		cache.WithBlock(),
	)
	s := storage.New(
		&dbConf{},
		storage.WithCache(c),
	)

	var (
		wg          sync.WaitGroup
		concurrency = 10
		id          = 238
		cacheKeyFmt = "content-id-%d"
		ctx         = storage.CtxWithCacheKey(context.Background(), fmt.Sprintf(cacheKeyFmt, id))
	)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()
			var content string
			_ = s.Find(ctx, &content, sqlFind, id)
			// require.Nil(t, err)
			log.Info("[%d] find content: %s, %p", i, content, &content)
		}(i)
	}
	wg.Wait()
	var content string
	_ = s.Find(ctx, &content, sqlFind, id)
	// require.Nil(t, err)
	log.Info("[main] find content: %s, %p", content, &content)
}

func TestStorage_Find(t *testing.T) {
	s := storage.New(
		&dbConf{},
	)
	var res string
	err := s.Find(context.Background(), &res, sqlFind, -1)
	assert.Equal(t, err, sql.ErrNoRows)
	err = s.Find(context.Background(), &res, sqlFind, 238)
	assert.Nil(t, err)
}

func TestStorage_Create_Find_Update_Delete(t *testing.T) {
	s := storage.New(
		&dbConf{},
		storage.WithCache(mem.New()),
	)
	var (
		err         error
		initContent = "chouchou"
		content     string
		id          int64
		ctx         = context.Background()
		cacheFmt    = "content-id-%d"
	)
	res, err := s.Create(ctx, sqlInsert, initContent)
	require.Nil(t, err)
	log.Info("create content: %s", initContent)

	id, _ = res.LastInsertId()
	ctx = storage.CtxWithCacheKey(ctx, fmt.Sprintf(cacheFmt, id))
	for i := 0; i < 2; i++ {
		err := s.Find(ctx, &content, sqlFind, id)
		require.Nil(t, err)
		log.Info("find content: %s", content)
	}
	err = s.Update(ctx, sqlUpdate, "1 "+content, id)
	require.Nil(t, err)
	log.Info("content updated")
	err = s.Find(ctx, &content, sqlFind, id)
	require.Nil(t, err)
	log.Info("find content: %s", content)
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
