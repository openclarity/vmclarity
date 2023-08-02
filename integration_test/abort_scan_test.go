package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/integration_test/helpers"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"time"
)

var _ = Describe("Aborting a scan", func() {

	Context("which is running", func() {
		It("should stop successfully", func(ctx SpecContext) {
			apiScanConfig, err := client.PostScanConfig(ctx, helpers.GetDefaultScanConfig())
			Expect(err).NotTo(HaveOccurred())

			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			odataFilter := fmt.Sprintf(
				"scanConfig/id eq '%s' and state ne '%s' and state ne '%s'",
				*apiScanConfig.Id,
				models.ScanStateDone,
				models.ScanStateFailed,
			)
			params := models.GetScansParams{
				Filter: &odataFilter,
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			err = client.PatchScan(ctx, *(*scans.Items)[0].Id, &models.Scan{
				State: utils.PointerTo(models.ScanStateAborted),
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
