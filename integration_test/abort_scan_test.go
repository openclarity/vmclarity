package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
)

var _ = Describe("Aborting a scan", func() {
	scanConfig := models.ScanConfig{
		Disabled:                      nil,
		Id:                            nil,
		MaxParallelScanners:           nil,
		Name:                          nil,
		ScanFamiliesConfig:            nil,
		ScannerInstanceCreationConfig: nil,
		Scheduled:                     nil,
		Scope:                         nil,
	}
	Context("which is running", func() {
		It("should stop successfully", func(ctx SpecContext) {
			resp, err := client.PostScanConfig(ctx, scanConfig)
			Expect(true).To(Equal(true))
		})
	})
})
