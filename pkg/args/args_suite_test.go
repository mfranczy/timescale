package args_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArgs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Args Suite")
}
