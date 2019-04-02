package collector

import (
	"fmt"
	"time"

	"github.com/mfranczy/timescale/pkg/container"
)

// Result represents processed data from worker pool
type Result struct {
	minT     int64
	maxT     int64
	medT     int64
	sumT     int64
	avgT     int64
	qSuccess int
	processT time.Duration
}

// PrintData shows processed data
func (r *Result) PrintData(jobsNum int) {
	fmt.Println("Total number of queries: ", jobsNum)
	fmt.Println("Number of succeeded queries: ", r.qSuccess)
	fmt.Println("Number of failed queries: ", jobsNum-r.qSuccess)
	fmt.Println("Total processing time: ", time.Duration(r.processT))
	fmt.Println("Maximum query time: ", time.Duration(r.maxT))
	fmt.Println("Minimum query time: ", time.Duration(r.minT))
	fmt.Println("Average query time: ", time.Duration(r.avgT))
	fmt.Println("Median query time: ", time.Duration(r.medT))
}

// Process data
func Process(dataCh <-chan int64) (rCh chan Result) {
	rCh = make(chan Result)
	result := Result{}
	maxHeap := container.NewMaxHeap()
	minHeap := container.NewMinHeap()

	go func() {
		start := time.Now()
		for val := range dataCh {
			if result.qSuccess == 0 {
				result.maxT = val
				result.minT = val
			} else {
				if val > result.maxT {
					result.maxT = val
				}
				if val < result.minT {
					result.minT = val
				}
			}
			result.qSuccess++
			result.sumT += val
			calculateMedian(maxHeap, minHeap, val)
		}
		result.medT = getMedian(maxHeap, minHeap)
		if result.qSuccess > 0 {
			result.avgT = result.sumT / int64(result.qSuccess)
		}
		result.processT = time.Since(start)
		rCh <- result
	}()

	return
}
