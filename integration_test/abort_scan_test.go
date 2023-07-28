package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
	"time"
)

var _ = Describe("Aborting a scan", func() {

	scope := "contains(assetInfo.labels, '{\"key\":\"scanconfig\",\"value\":\"test\"}')"
	newScanConfig := models.ScanConfig{
		Name: utils.PointerTo("Scan Config"),
		ScanTemplate: &models.ScanTemplate{
			AssetScanTemplate: &models.AssetScanTemplate{
				ScanFamiliesConfig: &models.ScanFamiliesConfig{
					Exploits: &models.ExploitsConfig{
						Enabled: utils.PointerTo(true),
					},
					Sbom: &models.SBOMConfig{
						Enabled: utils.PointerTo(true),
					},
					Vulnerabilities: &models.VulnerabilitiesConfig{
						Enabled: utils.PointerTo(true),
					},
				},
			},
			Scope: &scope,
		},
		Scheduled: &models.RuntimeScheduleScanConfig{
			CronLine: utils.PointerTo("0 */4 * * *"),
			OperationTime: utils.PointerTo(
				time.Date(2023, 1, 20, 15, 46, 18, 0, time.UTC),
			),
		},
	}

	Context("which is running", func() {
		It("should stop successfully", func(ctx SpecContext) {
			apiScanConfig, err := client.PostScanConfig(ctx, newScanConfig)
			Expect(err).NotTo(HaveOccurred())

			updateScanConfig := models.ScanConfig{
				Name: apiScanConfig.Name,
				ScanTemplate: &models.ScanTemplate{
					AssetScanTemplate: &models.AssetScanTemplate{
						ScanFamiliesConfig: apiScanConfig.ScanTemplate.AssetScanTemplate.ScanFamiliesConfig,
					},
					MaxParallelScanners: apiScanConfig.ScanTemplate.MaxParallelScanners,
					Scope:               apiScanConfig.ScanTemplate.Scope,
				},
				Scheduled: &models.RuntimeScheduleScanConfig{
					CronLine:      utils.PointerTo("0 */4 * * *"),
					OperationTime: utils.PointerTo(time.Now()),
				},
			}
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, &updateScanConfig)
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
