package container

type heap struct {
	eles []interface{}
	cmp  func(a, b interface{}) bool
}

func (h *heap) Insert(v interface{}) {
	h.eles = append(h.eles, v)
	h.up(len(h.eles) - 1)
}

func (h *heap) Remove() interface{} {
	v := h.eles[0]
	h.eles[0], h.eles[len(h.eles)-1] = h.eles[len(h.eles)-1], h.eles[0]
	h.eles = h.eles[:len(h.eles)-1]
	h.down(0)
	return v
}

func (h *heap) Top() interface{} {
	return h.eles[0]
}

func (h *heap) init() {
	for i := len(h.eles)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

func (h *heap) down(i int) {
	for {
		// find biggest child
		child := 2*i + 1
		if child > len(h.eles)-1 {
			break
		}
		if child+1 < len(h.eles) && h.cmp(h.eles[child+1], h.eles[child]) {
			child++
		}
		// swap if smaller than child
		if h.cmp(h.eles[i], h.eles[child]) {
			break
		}
		h.eles[i], h.eles[child] = h.eles[child], h.eles[i]
		i = child
	}
}

func (h *heap) up(i int) {
	for {
		// find parent
		if i == 0 {
			break
		}
		parent := (i - 1) / 2
		// swap if bigger than parent
		if h.cmp(h.eles[parent], h.eles[i]) {
			break
		}
		h.eles[i], h.eles[parent] = h.eles[parent], h.eles[i]
		i = parent
	}
}
