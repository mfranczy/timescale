package container

// Heap represents a heap of int64
type Heap []int64

// Len of heap
func (h Heap) Len() int { return len(h) }

// Swap two elements on heap
func (h Heap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push a new element to heap
func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

// Pop an element from heap
func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[0 : n-1]
	return val
}
