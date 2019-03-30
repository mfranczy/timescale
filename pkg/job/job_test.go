package job

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createCSVFile(dir, name string, data [][]string) (string, error) {
	fullPath := path.Join(dir, name)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

var _ = Describe("jobs", func() {
	var dir string
	var err error

	BeforeEach(func() {
		dir, err = ioutil.TempDir("/tmp", "timescale")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Context("with initializing new jobs", func() {

		It("should return an error when csv file does not exist", func() {
			jobs, err := New("file/not/found")
			Expect(err).To(HaveOccurred())
			Expect(jobs).To(BeNil())
		})

		It("should not return an error when csv file exists", func() {
			filePath, err := createCSVFile(dir, "test.csv", [][]string{{"hostname", "start_time", "end_time"}})
			Expect(err).NotTo(HaveOccurred())

			jobs, err := New(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobs).NotTo(BeNil())
			Expect(jobs.Num).To(Equal(0))

		})
	})

	Context("with csv dates validation", func() {

		It("should not return an error when data are correct", func() {
			row := []string{"host_test", "2017-01-02 18:21:03", "2017-01-02 19:21:03"}
			err := validateDates(row)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not have start_date higher than end_date", func() {
			row := []string{"host_test", "2017-01-04 18:21:03", "2017-01-02 19:21:03"}
			err := validateDates(row)
			Expect(err).To(HaveOccurred())
		})

		It("should validate date format of start_date", func() {
			row := []string{"host_test", "17-01-04 18:21:03", "2017-01-02 19:21:03"}
			err := validateDates(row)
			Expect(err).To(HaveOccurred())

			row = []string{"host_test", "test test test", "2017-01-02 19:21:03"}
			err = validateDates(row)
			Expect(err).To(HaveOccurred())

			row = []string{"host_test", "2017-01-02", "2017-01-02 19:21:03"}
			err = validateDates(row)
			Expect(err).To(HaveOccurred())
		})

		It("should validate date format of end_date", func() {
			row := []string{"host_test", "2017-01-02 19:21:03", "17-01-04 18:21:03"}
			err := validateDates(row)
			Expect(err).To(HaveOccurred())

			row = []string{"host_test", "2017-01-02 19:21:03", "test test test"}
			err = validateDates(row)
			Expect(err).To(HaveOccurred())

			row = []string{"host_test", "2017-01-02 19:21:03", "2017-01-02"}
			err = validateDates(row)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("with streaming jobs", func() {

		It("should return 0 jobs when there is nothing to provide to data channel", func(done Done) {
			filePath, err := createCSVFile(dir, "test.csv", [][]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).NotTo(Equal(""))

			jobs, err := New(filePath)
			Expect(err).NotTo(HaveOccurred())

			dataCh := jobs.Stream()
			job := <-dataCh
			Expect(job).To(Equal(Job{}))
			Expect(jobs.Num).To(Equal(0))

			close(done)
		}, 5)

		It("should stream data when csv file is not empty", func(done Done) {
			expectedJobs := []Job{
				{"host_test1", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
				{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"},
				{"host_test3", "2017-01-04 20:21:03", "2017-01-08 19:21:03"},
			}
			csvData := [][]string{
				{"hostname", "start_time", "end_time"},
				{"host_test1", "2017-01-02 18:21:03", "2017-01-02 19:21:03"},
				{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"},
				{"host_test3", "2017-01-04 20:21:03", "2017-01-08 19:21:03"},
			}
			filePath, err := createCSVFile(dir, "test.csv", csvData)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).NotTo(Equal(""))

			jobs, err := New(filePath)
			Expect(err).NotTo(HaveOccurred())
			dataCh := jobs.Stream()

			i := 0
			for d := range dataCh {
				Expect(d).To(Equal(expectedJobs[i]))
				i++
			}
			Expect(jobs.Num).To(Equal(len(expectedJobs)))
			Expect(i).To(Equal(len(expectedJobs)))

			close(done)
		}, 5)

		It("should not crash when csv file has corrupted data", func(done Done) {
			csvData := [][]string{
				{"hostname", "start_time", "end_time"},
				{"host_test1", "2017-10-02 18:21:03", "2017-01-02 19:21:03"},
				{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"},
				{"host_test3", "2017-01-04 20:21:03", "test test test"},
			}
			filePath, err := createCSVFile(dir, "test.csv", csvData)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).NotTo(Equal(""))

			jobs, err := New(filePath)
			Expect(err).NotTo(HaveOccurred())
			dataCh := jobs.Stream()

			i := 0
			for d := range dataCh {
				Expect(d).To(Equal(Job{"host_test2", "2017-01-03 19:21:03", "2017-01-05 21:21:03"}))
				i++
			}

			Expect(jobs.Num).To(Equal(len(csvData) - 1))
			Expect(i).To(Equal(1))

			close(done)
		}, 5)

	})
})
