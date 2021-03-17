package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestNewHeap(t *testing.T) {
	h := NewHeap([]interface{}{1, 2, 3}, func(a, b interface{}) bool {
		return a.(int) > b.(int)
	})
	assert.ElementsMatch(t, []interface{}{3, 2, 1}, h.(*heap).eles)
}

func TestHeap_Insert(t *testing.T) {
	h := NewHeap([]interface{}{}, func(a, b interface{}) bool {
		return len(a.(string)) > len(b.(string))
	})
	h.Insert("a")
	h.Insert("aa")
	h.Insert("aaa")
	assert.ElementsMatch(t, []interface{}{"aaa", "aa", "a"}, h.(*heap).eles)
}

func TestHeap_Remove(t *testing.T) {
	h := NewHeap([]interface{}{1, 2, 3, 4}, func(a, b interface{}) bool {
		return a.(int) > b.(int)
	})
	assert.ElementsMatch(t, []interface{}{4, 2, 3, 1}, h.(*heap).eles)
	assert.Equal(t, 4, h.Remove())
	assert.ElementsMatch(t, []interface{}{3, 2, 1}, h.(*heap).eles)
}
