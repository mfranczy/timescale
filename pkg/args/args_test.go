package args

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parse args", func() {

	Context("with specified arguments", func() {

		argsBackup := os.Args

		AfterEach(func() {
			os.Args = argsBackup
		})

		It("should not return an error when provided args are correct", func() {
			os.Args = append(os.Args, "-csv-file", "test.csv", "-workers-num", "4")
			a, err := Parse()

			Expect(err).NotTo(HaveOccurred())
			Expect(a.CSVFile).To(Equal("test.csv"))
			Expect(a.WorkersNum).To(Equal(4))
		})
	})
})
