package pool

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/mfranczy/timescale/pkg/db"
	"github.com/mfranczy/timescale/pkg/job"
)

var _ = Describe("worker", func() {

	Context("with db stmt", func() {

		It("should send an execution time", func(done Done) {
			_db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			defer _db.Close()

			mock.ExpectPrepare("SELECT (.*) FROM cpu_usage (.*)")
			mock.ExpectExec("SELECT (.*) FROM cpu_usage (.*)").WithArgs("host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03").
				WillReturnResult(sqlmock.NewResult(1, 1))

			stmt, err := db.GetCPUQuery(_db)
			Expect(err).NotTo(HaveOccurred())

			in := make(chan job.Job, 1)
			out := make(chan int64)
			finished := make(chan bool)

			worker := newWorker(1, in, out, finished)
			go worker.Exec(stmt)

			By("adding jobs to in channel")
			in <- job.Job{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"}
			close(in)

			By("reading jobs from out channel")
			execTime := <-out
			Expect(execTime).Should(BeNumerically(">", 0))

			By("waiting for worker to finish")
			Expect(<-finished).To(Equal(true))

			close(done)
		}, 5)

		It("should not send an execution time when error occured", func(done Done) {
			_db, mock, err := sqlmock.New()
			Expect(err).NotTo(HaveOccurred())
			defer _db.Close()

			mock.ExpectPrepare("SELECT (.*) FROM cpu_usage (.*)")
			mock.ExpectExec("SELECT (.*) FROM cpu_usage (.*)").WithArgs("host_test2", "invalid_date_format", "2017-01-05 21:21:03").
				WillReturnError(errors.New("Error"))

			stmt, err := db.GetCPUQuery(_db)
			Expect(err).NotTo(HaveOccurred())

			in := make(chan job.Job, 1)
			out := make(chan int64)
			finished := make(chan bool)

			worker := newWorker(1, in, out, finished)
			go worker.Exec(stmt)

			By("adding jobs to in channel")
			in <- job.Job{"host_test2", "invalid_date_format", "2017-01-05 21:21:03"}
			close(in)

			By("waiting for worker to finish")
			Expect(<-finished).To(Equal(true))

			close(done)
		}, 5)

	})
})
