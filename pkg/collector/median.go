package collector

import (
	"container/heap"

	"timescale/pkg/container"
)

func calculateMedian(maxHeap *container.MaxHeap, minHeap *container.MinHeap, val int64) {
	addElement(maxHeap, minHeap, val)
	balanceHeaps(maxHeap, minHeap)
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
		heap.Push(minHeap, heap.Pop(maxHeap))
		return
	}
	if minHeap.Len()-maxHeap.Len() >= 2 {
		heap.Push(maxHeap, heap.Pop(minHeap))
	}
}

func getMedian(maxHeap *container.MaxHeap, minHeap *container.MinHeap) int64 {
	if maxHeap.Len() == minHeap.Len() {
		return (heap.Pop(maxHeap).(int64) + heap.Pop(minHeap).(int64)) / int64(2)
	} else if maxHeap.Len() > minHeap.Len() {
		return heap.Pop(maxHeap).(int64)
	} else {
		return heap.Pop(minHeap).(int64)
	}
}
