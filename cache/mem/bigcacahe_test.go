package mem_test

import (
	"github.com/stretchr/testify/require"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/cache/mem"
	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.Init(
		log.WithWriters(&writer.Console{Formatter: &formatter.Text{}}),
		log.WithLevel(log.DebugLevel),
	)
	exitVal := m.Run()
	log.Stop()
	os.Exit(exitVal)
}

func TestNewBigCache(t *testing.T) {
	c := mem.NewBigCache(500, 3 * time.Second)
	_ = c.Set("k1", []byte("v1"))
	var v []byte
	err := c.Get("k1", &v)
	require.Nil(t, err)
	time.Sleep(5 *time.Second)
	err = c.Get("k1", &v)
	require.Error(t, cache.ErrNotFound)
}
