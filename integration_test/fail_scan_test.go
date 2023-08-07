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

var _ = Describe("Detecting scan failures", func() {

	Context("when a scan stops without assets to scan", func() {
		It("should detect failure reason successfully", func(ctx SpecContext) {
			By("applying a scan configuration with not existing label")
			apiScanConfig, err := client.PostScanConfig(
				ctx,
				helpers.GetCustomScanConfig(
					"contains(assetInfo.tags, '{\"key\":\"notexisting\",\"value\":\"label\"}')",
					1200,
				))
			Expect(err).NotTo(HaveOccurred())

			By("updating a scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan state changes to failed with nothing to scan as state reason")
			params := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s' and stateReason eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
					models.ScanRelationshipStateReasonNothingToScan,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())
		})
	})

	Context("when a scan stops with timeout", func() {
		It("should detect failure reason successfully", func(ctx SpecContext) {
			By("applying a scan configuration with short timeout")
			apiScanConfig, err := client.PostScanConfig(
				ctx,
				helpers.GetCustomScanConfig(
					"contains(assetInfo.labels, '{\"key\":\"scanconfig\",\"value\":\"test\"}')",
					2,
				))
			Expect(err).NotTo(HaveOccurred())

			By("updating a scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan starts")
			params := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state ne '%s' and state ne '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
					models.ScanStateFailed,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			By("waiting until scan state changes to failed with timed out as state reason")
			params = models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s' and stateReason eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateFailed,
					models.ScanStateReasonTimedOut,
				)),
			}
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())
		})
	})
})
