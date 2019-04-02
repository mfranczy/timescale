package pool

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mfranczy/timescale/pkg/job"
)

var _ = Describe("pool", func() {

	Context("with jobs assignement", func() {

		It("should send a job with the same hostname to the same worker as before", func(done Done) {
			expectedWorkersJobs := [][]job.Job{
				{
					// jobs with worker id = 0
					{"host_test1", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
				},
				{
					// jobs with worker id = 1
					{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"},
					{"host_test2", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
					{"host_test2", "2017-01-04 19:21:03", "2017-01-05 21:21:03"},
				},
				{
					// jobs with worker id = 2
					{"host_test3", "2017-01-04 20:21:03", "2017-01-08 19:21:03"},
					{"host_test3", "2017-01-05 20:21:03", "2017-01-08 19:21:03"},
				},
			}
			jobs := []job.Job{
				{"host_test1", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
				{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"},
				{"host_test3", "2017-01-04 20:21:03", "2017-01-08 19:21:03"},
				{"host_test2", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
				{"host_test2", "2017-01-04 19:21:03", "2017-01-05 21:21:03"},
				{"host_test3", "2017-01-05 20:21:03", "2017-01-08 19:21:03"},
			}
			jobsCh := make(chan job.Job, len(jobs))
			workersNum := 3
			p := New(workersNum)
			Expect(p.wNum).To(Equal(3))

			for _, j := range jobs {
				jobsCh <- j
			}
			close(jobsCh)
			p.dispatch(jobsCh)

			By("Checking worker with id 0")
			Expect(p.hostWorker["host_test1"]).To(Equal(0))
			Expect(len(p.inCh[0])).To(Equal(1))

			By("Checking worker with id 1")
			Expect(p.hostWorker["host_test2"]).To(Equal(1))
			Expect(len(p.inCh[1])).To(Equal(3))

			By("Checking worker with id 2")
			Expect(p.hostWorker["host_test3"]).To(Equal(2))
			Expect(len(p.inCh[2])).To(Equal(2))

			for i, w := range p.w {
				n := 0
				for j := range w.in {
					Expect(j).To(Equal(expectedWorkersJobs[i][n]))
					n++
				}
			}

			close(done)
		}, 5)

	})
})
