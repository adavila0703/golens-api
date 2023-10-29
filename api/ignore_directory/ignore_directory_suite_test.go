package ignore_directory_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIgnoreDirectory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IgnoreDirectory Suite")
}
