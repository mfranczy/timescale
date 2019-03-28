package collector

import (
	"container/heap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"timescale/pkg/container"
)

func addHeapElements(inf interface{}, val []int64) {
	for _, v := range val {
		switch inf.(type) {
		case *container.MaxHeap:
			heap.Push(inf.(*container.MaxHeap), v)
		case *container.MinHeap:
			heap.Push(inf.(*container.MinHeap), v)
		}
	}
}

var _ = Describe("calculate median", func() {

	Context("with data chain", func() {
		It("should return median with even length of data", func() {
			maxHeap := container.NewMaxHeap()
			minHeap := container.NewMinHeap()
			data := []int64{1, 2, 3, 4, 5, 6}

			for _, val := range data {
				calculateMedian(maxHeap, minHeap, val)
			}
			median := getMedian(maxHeap, minHeap)

			// division should be changed to float64 to get a real value
			Expect(median).To(Equal(int64(3)))
		})

		It("should return median with odd length of data", func() {
			maxHeap := container.NewMaxHeap()
			minHeap := container.NewMinHeap()
			data := []int64{1, 2, 3, 4, 5}

			for _, val := range data {
				calculateMedian(maxHeap, minHeap, val)
			}
			median := getMedian(maxHeap, minHeap)

			Expect(median).To(Equal(int64(3)))
		})
	})

	Context("with single values", func() {

		Context("with adding a new element to heap", func() {
			It("should add first element to max heap", func() {
				expectedMaxHeap := container.NewMaxHeap()
				addHeapElements(expectedMaxHeap, []int64{1})

				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addElement(maxHeap, minHeap, int64(1))
				Expect(maxHeap).To(Equal(expectedMaxHeap))
				Expect(minHeap.Len()).To(Equal(0))
			})

			It("should add smaller element to min heap", func() {
				expectedMaxHeap := container.NewMaxHeap()
				expectedMinHeap := container.NewMinHeap()
				addHeapElements(expectedMaxHeap, []int64{2, 3})
				addHeapElements(expectedMinHeap, []int64{4})

				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{3})
				addHeapElements(minHeap, []int64{4})

				addElement(maxHeap, minHeap, int64(2))
				Expect(maxHeap).To(Equal(expectedMaxHeap))
				Expect(minHeap).To(Equal(expectedMinHeap))
			})

			It("should add bigger element to max heap", func() {
				expectedMaxHeap := container.NewMaxHeap()
				expectedMinHeap := container.NewMinHeap()
				addHeapElements(expectedMaxHeap, []int64{3})
				addHeapElements(expectedMinHeap, []int64{4, 5})

				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{3})
				addHeapElements(minHeap, []int64{4})

				addElement(maxHeap, minHeap, int64(5))
				Expect(maxHeap).To(Equal(expectedMaxHeap))
				Expect(minHeap).To(Equal(expectedMinHeap))
			})
		})

		Context("with heap balancing", func() {

			By("adding expected elements to max and min heaps")
			expectedMaxHeap := container.NewMaxHeap()
			expectedMinHeap := container.NewMinHeap()
			addHeapElements(expectedMaxHeap, []int64{1, 2})
			addHeapElements(expectedMinHeap, []int64{3, 4})

			It("should pop element from min heap when max reached margin", func() {
				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{1})
				addHeapElements(minHeap, []int64{2, 3, 4})

				balanceHeaps(maxHeap, minHeap)
				Expect(maxHeap).To(Equal(expectedMaxHeap))
				Expect(minHeap).To(Equal(expectedMinHeap))
			})

			It("should pop element from max heap when min reached margin", func() {
				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{1, 2, 3})
				addHeapElements(minHeap, []int64{4})

				balanceHeaps(maxHeap, minHeap)
				Expect(maxHeap).To(Equal(expectedMaxHeap))
				Expect(minHeap).To(Equal(expectedMinHeap))
			})
		})

		Context("with getting median", func() {

			It("should return value from max heap", func() {
				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{3, 3})
				addHeapElements(minHeap, []int64{4})

				median := getMedian(maxHeap, minHeap)
				Expect(median).To(Equal(int64(3)))
			})

			It("should return value from min heap", func() {
				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{3})
				addHeapElements(minHeap, []int64{4, 5})

				median := getMedian(maxHeap, minHeap)
				Expect(median).To(Equal(int64(4)))
			})

			It("should return value divided from max and min heaps", func() {
				maxHeap := container.NewMaxHeap()
				minHeap := container.NewMinHeap()

				addHeapElements(maxHeap, []int64{3})
				addHeapElements(minHeap, []int64{5})

				median := getMedian(maxHeap, minHeap)
				Expect(median).To(Equal(int64(4)))
			})
		})
	})
})
