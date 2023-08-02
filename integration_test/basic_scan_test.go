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

var _ = Describe("Running a basic scan", func() {

	Context("which scans a docker container", func() {
		It("should finish successfully", func(ctx SpecContext) {

			// Check if asset discovered
			assetsParams := models.GetAssetsParams{
				Filter: utils.PointerTo(fmt.Sprintf("assetInfo.containerName eq '/vmclarity-ubuntu-1'")),
			}
			Eventually(func() bool {
				assets, err := client.GetAssets(ctx, assetsParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*assets.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			apiScanConfig, err := client.PostScanConfig(ctx, helpers.GetDefaultScanConfig())
			Expect(err).NotTo(HaveOccurred())

			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			// Check if scan is running
			scanParams := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state ne '%s' and state ne '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
					models.ScanStateFailed,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, scanParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			// Check if scan is finished
			scanParams = models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
				)),
			}
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, scanParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*360, time.Second).Should(BeTrue())
		})
	})
})
