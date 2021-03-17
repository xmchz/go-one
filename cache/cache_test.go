package cache_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xmchz/go-one/cache"
	"github.com/xmchz/go-one/cache/mem"
	"github.com/xmchz/go-one/log"
	"github.com/xmchz/go-one/log/formatter"
	"github.com/xmchz/go-one/log/writer"
	metricMem "github.com/xmchz/go-one/metric/mem"
	"os"
	"sync"
	"testing"
)

var source = func(v interface{}) error {
	v = 1
	log.Debug("from source")
	return nil
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

func TestWithBlock(t *testing.T) {
	c := cache.New(
		mem.New(),
		cache.WithBlock(),
	)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			var v int
			err := c.Take(&v, "k", source)
			assert.Nil(t, err)
		}()
	}
	wg.Wait()
}

func TestWithMetric(t *testing.T) {
	counter := metricMem.NewCounter()
	c := cache.New(
		mem.New(),
		cache.WithMetric(
			counter,
		),
	)
	var v int
	err := c.Take(&v, "k", source)
	assert.Nil(t, err)
	err = c.Get("k", &v)
	assert.Nil(t, err)
	err = c.Get("k", &v)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(counter.Value()))
}
