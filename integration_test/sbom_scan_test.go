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

var _ = Describe("Running a basic scan (only SBOM)", func() {

	Context("which scans a docker container", func() {
		It("should finish successfully", func(ctx SpecContext) {

			By("waiting until test asset is found")
			assetsParams := models.GetAssetsParams{
				Filter: utils.PointerTo("assetInfo.containerName eq '/alpine-test'"),
				//TODO(paralta) Filter: utils.PointerTo("assetInfo/labels/any(l: l/Key eq 'scanconfig' and l/Value eq 'test')"),
			}
			Eventually(func() bool {
				assets, err := client.GetAssets(ctx, assetsParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*assets.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			By("applying a scan configuration")
			apiScanConfig, err := client.PostScanConfig(ctx, helpers.GetSBOMScanConfig())
			Expect(err).NotTo(HaveOccurred())

			By("updating a scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan starts")
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

			By("waiting until scan state changes to done")
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
			}, time.Second*120, time.Second).Should(BeTrue())
		})
	})
})
