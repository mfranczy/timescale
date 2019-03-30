package container

import (
	"container/heap"
)

// MaxHeap represents a max heap
type MaxHeap struct {
	Heap
}

// Less for max heap
func (m *MaxHeap) Less(i int, j int) bool {
	return m.Heap[i] > m.Heap[j]
}

// MinHeap represents a min heap
type MinHeap struct {
	Heap
}

// Less for min heap
func (m *MinHeap) Less(i int, j int) bool {
	return m.Heap[i] < m.Heap[j]
}

// NewMaxHeap initializes a new max heap
func NewMaxHeap() *MaxHeap {
	h := &MaxHeap{[]int64{}}
	heap.Init(h)
	return h
}

// NewMinHeap initializes a new min heap
func NewMinHeap() *MinHeap {
	h := &MinHeap{[]int64{}}
	heap.Init(h)
	return h
}

// AddHeapElements adds elements to max or min heap
func AddHeapElements(inf interface{}, val []int64) {
	for _, v := range val {
		switch inf.(type) {
		case *MaxHeap:
			heap.Push(inf.(*MaxHeap), v)
		case *MinHeap:
			heap.Push(inf.(*MinHeap), v)
		}
	}
}
