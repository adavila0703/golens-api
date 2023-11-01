package health_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"golens-api/api/health"
)

var _ = Describe("HealthCheck", func() {
	It("runs a health check", func() {
		res, err := health.HealthCheck(nil, nil, nil)
		resMessage := res.(*health.HealthCheckResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Message).To(Equal("Good!"))
	})
})
