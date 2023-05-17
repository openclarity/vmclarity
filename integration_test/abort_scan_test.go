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
	scope := models.AwsScanScope{
		AllRegions: utils.PointerTo(true),
		InstanceTagSelector: &[]models.Tag{
			{
				Key:   "ScanConfig",
				Value: "test",
			},
		},
		ObjectType: "AwsScanScope",
	}
	var scanScopeType models.ScanScopeType
	err := scanScopeType.FromAwsScanScope(scope)
	Expect(err).NotTo(HaveOccurred())

	newScanConfig := models.ScanConfig{
		Name: utils.PointerTo("Scan Config"),
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
		Scheduled: &models.RuntimeScheduleScanConfig{
			CronLine: utils.PointerTo("0 */4 * * *"),
			OperationTime: utils.PointerTo(
				time.Date(2023, 1, 20, 15, 46, 18, 0, time.UTC),
			),
		},
		Scope: &scanScopeType,
	}

	Context("which is running", func() {
		It("should stop successfully", func(ctx SpecContext) {
			apiScanConfig, err := client.PostScanConfig(ctx, newScanConfig)
			Expect(err).NotTo(HaveOccurred())

			updateScanConfig := models.ScanConfig{
				MaxParallelScanners: apiScanConfig.MaxParallelScanners,
				Name:                apiScanConfig.Name,
				ScanFamiliesConfig:  apiScanConfig.ScanFamiliesConfig,
				Scheduled: &models.RuntimeScheduleScanConfig{
					CronLine:      utils.PointerTo("0 */4 * * *"),
					OperationTime: utils.PointerTo(time.Now()),
				},
				Scope: apiScanConfig.Scope,
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
