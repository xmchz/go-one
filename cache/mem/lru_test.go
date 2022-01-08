package mem_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/cache/mem"
	"github.com/xmchz/go-one/log"
	logCore "github.com/xmchz/go-one/log/core"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
)

func TestMain(m *testing.M) {
	log.Init(
		logCore.WithWriters(writer.NewConsole(&formatter.Text{})),
		logCore.WithLevel(logCore.DebugLevel),
	)
	exitVal := m.Run()
	log.Stop()
	os.Exit(exitVal)
}

func TestNewLru(t *testing.T) {
	c := mem.New()
	err := c.Set("k1", []byte("v1"))
	require.Nil(t, err)
	log.Info("set k1 success")

	var v1 []byte
	err = c.Get("k1", &v1)
	require.Nil(t, err)
	log.Info("get k1: %#v", v1)

	var v2 []byte
	err = c.Get("k2", &v2)
	require.Error(t, err, cache.ErrNotFound)
	log.Info("get k1: %#v", v2)
}
