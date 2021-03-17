package container

func NewLru(capacity int) *Lru {
	head := &entity{}
	tail := &entity{}
	head.next = tail
	tail.pre = head
	return &Lru{
		m:        make(map[interface{}]*entity),
		capacity: capacity,
		head:     head,
		tail:     tail,
	}
}

type Lru struct {
	m        map[interface{}]*entity
	capacity int
	head     *entity
	tail     *entity
}

type entity struct {
	key  interface{}
	val  interface{}
	pre  *entity
	next *entity
}

func (l *Lru) Get(key interface{}) (interface{}, bool) {
	e, ok := l.m[key]
	if !ok {
		return nil, false
	}
	l.moveToHead(e)
	return e.val, true
}

func (l *Lru) Set(key interface{}, val interface{}) {
	e, ok := l.m[key]
	if !ok {
		if len(l.m) + 1 > l.capacity {
			eldest := l.tail.pre
			l.remove(eldest)
			delete(l.m, eldest.key)
		}
		e = &entity{
			key: key,
			val: val,
		}
		l.add(e)
		l.m[key] = e

	} else {
		e.val = val
		l.moveToHead(e)
	}
}

func (l *Lru) Del(key interface{}) {
	if e, ok := l.m[key]; ok {
		l.remove(e)
		delete(l.m, e.key)
	}
}

func (l *Lru) add(e *entity) {
	e.next = l.head.next
	l.head.next = e
	e.pre = l.head
	e.next.pre = e
}

func (l *Lru) remove(e *entity) {
	e.pre.next = e.next
	e.next.pre = e.pre
}

func (l *Lru) moveToHead(e *entity) {
	l.remove(e)
	l.add(e)
}
