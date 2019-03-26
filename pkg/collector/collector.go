package collector

import (
	"fmt"
	"time"

	"timescale/pkg/container"
)

// Result represents processed data from worker pool
type Result struct {
	minT     int64
	maxT     int64
	medT     int64
	sumT     int64
	qSuccess int64
	processT time.Duration
}

// PrintData shows processed data
func (r *Result) PrintData(jobsNum int) {
	avgT := int64(0)
	if r.qSuccess > 0 {
		avgT = r.sumT / r.qSuccess
	}

	fmt.Println("Total number of queries: ", jobsNum)
	fmt.Println("Number of succeeded queries: ", r.qSuccess)
	fmt.Println("Number of failed queries: ", int64(jobsNum)-r.qSuccess)
	fmt.Println("Total processing time: ", time.Duration(r.processT))
	fmt.Println("Maximum query time: ", time.Duration(r.maxT))
	fmt.Println("Minimum query time: ", time.Duration(r.minT))
	fmt.Println("Average query time: ", time.Duration(avgT))
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
			result.medT = getMedian(maxHeap, minHeap, val)
		}
		result.processT = time.Since(start)
		rCh <- result
	}()

	return
}
