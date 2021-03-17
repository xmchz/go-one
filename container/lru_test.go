package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLru(t *testing.T) {
	l := NewLru(3)
	l.Set(1, "a")
	l.Set(2, "b")
	l.Set(3, "c")
	l.Set(4, "d")
	_, ok := l.Get(1)
	assert.False(t, ok)
	v, _ := l.Get(2)
	assert.Equal(t, "b", v)
	l.Set(5, "e")
	_, ok = l.Get(3)
	assert.False(t, ok)
}
