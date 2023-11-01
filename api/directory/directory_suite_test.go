package directory_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDirectory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Directory Suite")
}
