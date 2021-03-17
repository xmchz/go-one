package container

import "strings"

func newTrie() *trie {
	return &trie{
		root: &trieNode{
			children: make(map[rune]*trieNode, 32),
		},
	}
}

type trie struct {
	root *trieNode
}

type trieNode struct {
	char     rune
	data     interface{}
	Depth    int
	parent   *trieNode
	children map[rune]*trieNode
	term     bool
}

func (p *trie) Add(key string, data interface{}) {
	cur := p.root
	for _, r := range []rune(strings.TrimSpace(key)) {
		n, ok := cur.children[r]
		if !ok {
			n = &trieNode{
				children: make(map[rune]*trieNode, 32),
			}
			n.Depth = cur.Depth + 1
			n.char = r
			cur.children[r] = n
		}
		cur = n
	}
	cur.term = true
	cur.data = data
}

func (p *trie) PrefixSearch(key string) []interface{} {
	n := p.find(key)
	if n == nil {
		return nil
	}
	ns := p.collect(n)
	res := make([]interface{}, len(ns))
	for i, n := range ns {
		res[i] = n.data
	}
	return res
}

func (p *trie) Replace(origin, replace string) (string, bool) {
	if p.root == nil {
		return "", false
	}

	var hit bool
	runes := []rune(strings.TrimSpace(origin))
	var left []rune
	cur := p.root
	start := 0
	for idx, r := range runes {
		n, ok := cur.children[r]
		if !ok {
			left = append(left, runes[start:idx+1]...)
			start = idx + 1
			cur = p.root
			continue
		}
		if n.term {
			hit = true
			left = append(left, []rune(replace)...)
			start = idx + 1
			cur = p.root
		} else {
			cur = n
		}
	}
	return string(left), hit
}

func (p *trie) find(key string) *trieNode {
	cur := p.root
	for _, r := range []rune(strings.TrimSpace(key)) {
		n, ok := cur.children[r]
		if !ok {
			return nil
		}
		cur = n
	}
	return cur
}

func (p *trie) collect(n *trieNode) []*trieNode {
	if n == nil {
		return nil
	}
	if n.term {
		return []*trieNode{n}
	}
	var res []*trieNode
	q := []*trieNode{n}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.term {
			res = append(res, cur)
		} else {
			for _, child := range cur.children {
				q = append(q, child)
			}
		}
	}
	return res
}
