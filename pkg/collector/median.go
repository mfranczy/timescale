package collector

import (
	"container/heap"
	"timescale/pkg/container"
)

func getMedian(maxHeap *container.MaxHeap, minHeap *container.MinHeap, val int64) int64 {
	addElement(maxHeap, minHeap, val)
	balanceHeaps(maxHeap, minHeap)
	return calcMedian(maxHeap, minHeap)
}

func addElement(maxHeap *container.MaxHeap, minHeap *container.MinHeap, val int64) {
	if maxHeap.Len() == 0 || val <= maxHeap.Heap[0] {
		heap.Push(maxHeap, val)
	} else {
		heap.Push(minHeap, val)
	}
}

func balanceHeaps(maxHeap *container.MaxHeap, minHeap *container.MinHeap) {
	if maxHeap.Len()-minHeap.Len() >= 2 {
		heap.Push(minHeap, maxHeap.Pop())
		return
	}
	if minHeap.Len()-maxHeap.Len() >= 2 {
		heap.Push(maxHeap, minHeap.Pop())
	}
}

func calcMedian(maxHeap *container.MaxHeap, minHeap *container.MinHeap) int64 {
	if maxHeap.Len() == minHeap.Len() {
		return (maxHeap.Pop().(int64) + minHeap.Pop().(int64)) / int64(2)
	} else if maxHeap.Len() > minHeap.Len() {
		return maxHeap.Pop().(int64)
	} else {
		return minHeap.Pop().(int64)
	}
}
