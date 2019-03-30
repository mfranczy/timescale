package container

import (
	"container/heap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("container", func() {

	Context("with max and min heaps", func() {

		It("should pop elements in the right order", func() {
			expectedHeapOrder := []int64{1, 2, 3, 4, 5}
			maxHeap := NewMaxHeap()
			minHeap := NewMinHeap()

			AddHeapElements(maxHeap, []int64{1, 3, 5, 4, 2})
			AddHeapElements(minHeap, []int64{1, 3, 5, 4, 2})

			for i, j := 0, 4; i < 5; i, j = i+1, j-1 {
				Expect(heap.Pop(maxHeap)).To(Equal(expectedHeapOrder[j]))
				Expect(heap.Pop(minHeap)).To(Equal(expectedHeapOrder[i]))
			}
		})
	})

})
