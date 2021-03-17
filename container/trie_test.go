package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_Replace(t *testing.T) {
	trie := NewTrie()
	trie.Add("黄色", "黄色")
	trie.Add("黄绿色", "黄绿色")
	trie.Add("蓝色", "蓝色")

	res1, ok := trie.Replace("黄色的灯泡是蓝色的", "***")
	res2 := trie.PrefixSearch("黄")
	assert.Equal(t, "***的灯泡是***的", res1)
	assert.Equal(t, true, ok)
	assert.ElementsMatch(t, []interface{}{"黄色", "黄绿色"}, res2)
}
