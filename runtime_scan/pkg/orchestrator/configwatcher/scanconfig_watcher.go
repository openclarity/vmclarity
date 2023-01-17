// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configwatcher

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
)

const (
	timeWindow = 5 * time.Minute
)

type ScanConfigWatcher struct {
	backendClient  *client.ClientWithResponses
	providerClient provider.Client
	scannerConfig  *_config.ScannerConfig
	cancelFn       context.CancelFunc
}

func CreateScanConfigWatcher(
	backendClient *client.ClientWithResponses,
	providerClient provider.Client,
	scannerConfig _config.ScannerConfig,
) *ScanConfigWatcher {
	return &ScanConfigWatcher{
		backendClient:  backendClient,
		providerClient: providerClient,
		scannerConfig:  &scannerConfig,
	}
}

func (scw *ScanConfigWatcher) SetCancelFn(cancel context.CancelFunc) {
	scw.cancelFn = cancel
}

func (scw *ScanConfigWatcher) getScanConfigs() (*models.ScanConfigs, error) {
	resp, err := scw.backendClient.GetScanConfigsWithResponse(context.TODO(), &models.GetScanConfigsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get a scan configs: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scan configs: empty body")
		}
		return resp.JSON200, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get scan configs. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get scan configs. status code=%v", resp.StatusCode())
	}
}

func (scw *ScanConfigWatcher) getScansByScanConfigIDAndOperationTime(scanConfigID string, operationTime time.Time) ([]models.Scan, error) {
	//odataFilter := fmt.Sprintf("scanConfigId eq '%s' and (endTime eq null or startTime gte '%s')", scanConfigID, operationTime.String())
	//params := &models.GetScansParams{
	//	Filter: &odataFilter,
	//}
	resp, err := scw.backendClient.GetScansWithResponse(context.TODO(), &models.GetScansParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get a scans with: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scans: empty body")
		}
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get scans. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get scans. status code=%v", resp.StatusCode())
	}
	// After Odata filters will be implemented on the backend the filter function can be removed
	return scw.filterScanConfigs(resp.JSON200, scanConfigID, operationTime), nil
}

func (scw *ScanConfigWatcher) getScanConfigsToScan() ([]models.ScanConfig, error) {
	scanConfigsToScan := make([]models.ScanConfig, 0)
	scanConfigs, err := scw.getScanConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to check new scan configs: %v", err)
	}

	now := time.Now()
	for _, scanConfig := range *scanConfigs.Items {
		// Check only the SingleScheduledScanConfigs at the moment
		operationTime, ok, err := getSingleScheduledScanConfigOperationTime(scanConfig.Scheduled)
		if err != nil {
			log.Errorf("Failed to check scan config type with id=%s: %v", *scanConfig.Id, err)
			continue
		}
		if !ok {
			continue
		}
		// ScanConfig skip to start because its within the window
		if operationTime.Sub(now).Abs() >= timeWindow {
			continue
		}

		// Need to check existing Scans to determine if we can create a Scan
		scans, err := scw.getScansByScanConfigIDAndOperationTime(*scanConfig.Id, operationTime)
		if err != nil {
			log.Errorf("Failed to get scan configs: %v", err)
			continue
		}
		if len(scans) != 0 {
			scanConfigsToScan = append(scanConfigsToScan, scanConfig)
		}
	}
	return scanConfigsToScan, nil
}

func (scw *ScanConfigWatcher) filterScanConfigs(scans *models.Scans, scanConfigID string, operationTime time.Time) []models.Scan {
	for _, scan := range *scans.Items {
		if *scan.ScanConfigId != scanConfigID {
			continue
		}
		if scan.EndTime == nil {
			// there is a running scan for this scanConfig
			continue
		}
		if scan.StartTime.After(operationTime) {
			// there is already a scan created for this scanConfig operation time
			continue
		}
		return []models.Scan{scan}
	}
	return []models.Scan{}
}

func getSingleScheduledScanConfigOperationTime(scheduleScanConfigType *models.RuntimeScheduleScanConfigType) (time.Time, bool, error) {
	scanConfig, err := scheduleScanConfigType.ValueByDiscriminator()
	if err != nil {
		return time.Time{}, false, fmt.Errorf("failed to determine scheduled scan config type: %v", err)
	}
	switch scanConfig.(type) {
	case models.SingleScheduleScanConfig:
		// nolint:forcetypeassert
		return scanConfig.(*models.SingleScheduleScanConfig).OperationTime, true, nil
	default:
		return time.Time{}, false, nil
	}
}

func (scw *ScanConfigWatcher) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(scw.scannerConfig.ScanConfigWatchInterval):
				scanConfigsToScan, err := scw.getScanConfigsToScan()
				if err != nil {
					log.Warnf("Failed to check scan configs: %v", err)
				}
				if len(scanConfigsToScan) > 0 {
					scw.scheduleNewScans(scanConfigsToScan)
				}
			case <-ctx.Done():
				log.Infof("Stop watching scan configs.")
				return
			}
		}
	}()
}

func (scw *ScanConfigWatcher) Stop() {
	scw.cancelFn()
}
