package container

type Ring interface {
	Insert(interface{}) bool
	Remove() (interface{}, bool)
}

type Heap interface {
	Insert(interface{})
	Remove() interface{}
	Top() interface{}
}

type Trie interface {
	Add(string, interface{})
	PrefixSearch(key string) []interface{}
	Replace(origin, replace string) (string, bool)
}

func NewHeap(arr []interface{}, cmp func(a, b interface{}) bool) Heap {
	h := &heap{
		eles: arr,
		cmp:  cmp,
	}
	h.init()
	return h
}

func NewRing(size int) Ring {
	return &ring{
		eles: make([]interface{}, size),
		tail: -1,
		head: -1,
	}
}

func NewTrie() Trie {
	return newTrie()
}

/*
New() 说明包中只有一个实现，包名更加具体，或作为lib的公开接口
NewXXX() 说明包中有一类实现，包名泛指一类功能：container service storage cache等
 */
