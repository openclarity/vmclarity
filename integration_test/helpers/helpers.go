package helpers

import (
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"time"
)

func GetDefaultScanConfig() models.ScanConfig {
	return models.ScanConfig{
		Name: utils.PointerTo(uuid.New().String()),
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
			Scope: utils.PointerTo("contains(assetInfo.labels, '{\"key\":\"scanconfig\",\"value\":\"test\"}')"),
		},
		Scheduled: &models.RuntimeScheduleScanConfig{
			CronLine: utils.PointerTo("0 */4 * * *"),
			OperationTime: utils.PointerTo(
				time.Date(2023, 1, 20, 15, 46, 18, 0, time.UTC),
			),
		},
	}
}

func UpdateScanConfigToStartNow(config *models.ScanConfig) *models.ScanConfig {
	return &models.ScanConfig{
		Name: config.Name,
		ScanTemplate: &models.ScanTemplate{
			AssetScanTemplate: &models.AssetScanTemplate{
				ScanFamiliesConfig: config.ScanTemplate.AssetScanTemplate.ScanFamiliesConfig,
			},
			MaxParallelScanners: config.ScanTemplate.MaxParallelScanners,
			Scope:               config.ScanTemplate.Scope,
		},
		Scheduled: &models.RuntimeScheduleScanConfig{
			CronLine:      utils.PointerTo("0 */4 * * *"),
			OperationTime: utils.PointerTo(time.Now()),
		},
	}
}
