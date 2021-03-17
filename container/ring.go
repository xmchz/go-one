package container

type ring struct {
	eles       []interface{}
	tail, head int
}

func (r *ring) Insert(v interface{}) bool {
	if r.full() {
		return false
	}
	if r.empty() {
		r.head = 0
		r.tail = 0
		r.eles[0] = v
		return true
	}

	r.tail++
	if r.tail == len(r.eles) {
		r.tail = 0
	}
	r.eles[r.tail] = v
	return true
}

func (r *ring) Remove() (interface{}, bool) {
	if r.empty() {
		return nil, false
	}
	v := r.eles[r.head]
	if r.head == r.tail {
		r.head = -1
		r.tail = -1
	} else {
		r.head++
		if r.head == len(r.eles) {
			r.head = 0
		}
	}
	return v, true
}

func (r *ring) full() bool {
	return r.tail-r.head+1 == len(r.eles) || r.head-r.tail == 1
}

func (r *ring) empty() bool {
	//return r.tail == r.head
	return r.tail == -1 && r.head == -1
}
