package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRing_Insert(t *testing.T) {
	r := NewRing(4)
	var res1 []bool
	for i := 1; i <= 5; i++ {
		res1 = append(res1, r.Insert(i))
	}
	var (
		res2 []interface{}
		res3 []bool
	)
	for i := 1; i <= 5; i++ {
		v, ok := r.Remove()
		res2 = append(res2, v)
		res3 = append(res3, ok)
	}
	assert.ElementsMatch(t, []bool{true, true, true, true, false}, res1)
	assert.ElementsMatch(t, []interface{}{1, 2, 3, 4, nil}, res2)
	assert.ElementsMatch(t, []interface{}{1, 2, 3, 4}, r.(*ring).eles)
	assert.ElementsMatch(t, []bool{true, true, true, true, false}, res3)
}
