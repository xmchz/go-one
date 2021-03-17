package mem_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xmchz/go-one/metric/mem"
	"testing"
)

func TestCounter(t *testing.T) {
	c := mem.NewCounter()
	c.Add(1)
	c.Add(1)
	assert.Equal(t, 2, int(c.Value()))
}
