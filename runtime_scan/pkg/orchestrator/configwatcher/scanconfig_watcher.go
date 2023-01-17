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
)

const (
	pageSize      = 100
	timeWindowMin = 5
)

type ScanConfigWatcher struct {
	stopChan                chan struct{}
	scanConfigsChan         chan []models.ScanConfig
	backendClient           *client.ClientWithResponses
	scanConfigWatchInterval time.Duration
}

func CreateScanConfigWatcher(scanConfigsChan chan []models.ScanConfig,
	backendClient *client.ClientWithResponses,
) *ScanConfigWatcher {
	return &ScanConfigWatcher{
		stopChan:        make(chan struct{}),
		scanConfigsChan: scanConfigsChan,
		backendClient:   backendClient,
	}
}

func (scw *ScanConfigWatcher) getScanConfigs() ([]models.ScanConfig, error) {
	params := &models.GetScanConfigsParams{
		Page:     1,
		PageSize: pageSize,
	}

	var scanConfigs []models.ScanConfig
	paginatedScanConfigs, err := scw.getScanConfigsByPage(params)
	if err != nil {
		return scanConfigs, fmt.Errorf("failed to get scan configs by page:%d, pageSize: %d, error:%v", params.Page, params.PageSize, err)
	}
	scanConfigs = append(scanConfigs, *paginatedScanConfigs.Items...)
	for {
		if paginatedScanConfigs.Total-(params.Page*params.PageSize) <= 0 {
			break
		}
		params.Page++
		paginatedScanConfigs, err = scw.getScanConfigsByPage(params)
		if err != nil {
			return scanConfigs, fmt.Errorf("failed to get scanconfigs by page:%d, pageSize: %d, error:%v", params.Page, params.PageSize, err)
		}
		scanConfigs = append(scanConfigs, *paginatedScanConfigs.Items...)
	}

	return scanConfigs, nil
}

func (scw *ScanConfigWatcher) getScanConfigsByPage(params *models.GetScanConfigsParams) (*models.ScanConfigs, error) {
	resp, err := scw.backendClient.GetScanConfigsWithResponse(context.TODO(), params)
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

func (scw *ScanConfigWatcher) getScansByScanConfigID(scanConfigID string) ([]models.Scan, error) {
	odataFilter := fmt.Sprintf("scanConfigId eq '%s'", scanConfigID)
	params := &models.GetScansParams{
		Filter:   &odataFilter,
		Page:     1,
		PageSize: pageSize,
	}

	var scans []models.Scan
	paginatedScans, err := scw.getScansByPage(params)
	if err != nil {
		return scans, fmt.Errorf("failed to get scans by page:%d, pageSize: %d, error:%v", params.Page, params.PageSize, err)
	}
	scans = append(scans, *paginatedScans.Items...)
	for {
		if paginatedScans.Total-(params.Page*params.PageSize) <= 0 {
			break
		}
		params.Page++
		paginatedScans, err = scw.getScansByPage(params)
		if err != nil {
			return scans, fmt.Errorf("failed to get scanconfigs by page:%d, pageSize: %d, error:%v", params.Page, params.PageSize, err)
		}
		scans = append(scans, *paginatedScans.Items...)
	}

	return scans, nil
}

func (scw *ScanConfigWatcher) getScansByPage(params *models.GetScansParams) (*models.Scans, error) {
	resp, err := scw.backendClient.GetScansWithResponse(context.TODO(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get a scans with params=%s: %v", *params.Filter, err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scans: empty body")
		}
		return resp.JSON200, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get scans. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get scans. status code=%v", resp.StatusCode())
	}
}

func (scw *ScanConfigWatcher) checkScanConfigs() ([]models.ScanConfig, error) {
	scanConfigsToScan := make([]models.ScanConfig, 0)
	scanConfigs, err := scw.getScanConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to check new scan configs: %v", err)
	}

	for _, scanConfig := range scanConfigs {
		// Check only the SingleScheduledScanConfigs at the moment
		operationTime, ok, err := getSingleScheduledScanConfigOperationTime(scanConfig.Scheduled)
		if err != nil {
			log.Errorf("Failed to check scan config type with id=%s: %v", *scanConfig.Id, err)
			continue
		}
		if !ok {
			continue
		}
		now := time.Now()
		// ScanConfig needs to start because its within the window
		if !(operationTime.Before(now.Add(timeWindowMin*time.Minute)) && operationTime.After(now.Add(-timeWindowMin*time.Minute))) {
			continue
		}

		// Need to check existing Scans to determine if we can create a Scan
		initiate, err := scw.shouldInitiateScanForScanConfig(*scanConfig.Id, operationTime)
		if err != nil {
			log.Errorf("Failed to determine whether scan should be initiated: %v", err)
			continue
		}
		if initiate {
			scanConfigsToScan = append(scanConfigsToScan, scanConfig)
		}
	}
	return scanConfigsToScan, nil
}

func (scw *ScanConfigWatcher) shouldInitiateScanForScanConfig(scanConfigID string, operationTime time.Time) (bool, error) {
	scans, err := scw.getScansByScanConfigID(scanConfigID)
	if err != nil {
		return false, fmt.Errorf("failed to get scans for scan config ID=%s: %v", scanConfigID, err)
	}
	for _, scan := range scans {
		if scan.EndTime == nil {
			// there is a running scan for this scanConfig
			return false, nil
		}
		if scan.StartTime.After(operationTime) {
			// there is already a scan created for this scanConfig operation
			return false, nil
		}
	}
	return true, nil
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

func (scw *ScanConfigWatcher) Start() {
	// Clear
	close(scw.stopChan)
	scw.stopChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(scw.scanConfigWatchInterval):
				scanConfigsToScan, err := scw.checkScanConfigs()
				if err != nil {
					log.Warnf("Failed to check scan configs: %v", err)
				}
				if len(scanConfigsToScan) > 0 {
					scw.scanConfigsChan <- scanConfigsToScan
				}
			case <-scw.stopChan:
				log.Infof("Stop watching scan configs.")
				return
			}
		}
	}()
}

func (scw *ScanConfigWatcher) Stop() {
	close(scw.stopChan)
}
