package collector

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("collector", func() {
	buffer := 5

	Context("when data channel is empty", func() {
		It("should return empty result", func() {
			dataCh := make(chan int64, buffer)
			resCh := Process(dataCh)
			close(resCh)
			res := <-resCh

			Expect(res.maxT).To(Equal(int64(0)))
			Expect(res.minT).To(Equal(int64(0)))
			Expect(res.sumT).To(Equal(int64(0)))
			Expect(res.medT).To(Equal(int64(0)))
			Expect(res.avgT).To(Equal(int64(0)))
			Expect(res.qSuccess).To(Equal(0))
		})
	})

	Context("when data channel contains output from workers", func() {
		It("should calculate avg, max, median, # of jobs, processed time", func() {
			dataCh := make(chan int64, buffer)
			data := []int64{1, 2, 3, 4, 5}

			for _, val := range data {
				dataCh <- val
			}
			close(dataCh)

			resCh := Process(dataCh)
			res := <-resCh

			Expect(res.maxT).To(Equal(int64(5)))
			Expect(res.minT).To(Equal(int64(1)))
			Expect(res.sumT).To(Equal(int64(15)))
			Expect(res.medT).To(Equal(int64(3)))
			Expect(res.avgT).To(Equal(int64(3)))
			Expect(res.qSuccess).To(Equal(buffer))
		})
	})
})
