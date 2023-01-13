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

package scanconfig_watcher

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
)

type ScanConfigWatcher struct {
	existingScanConfigMap   map[string]models.ScanConfig
	stopChan                chan struct{}
	scanConfigChan          chan *map[string]models.ScanConfig
	backendClient           *client.ClientWithResponses
	scanConfigWatchInterval time.Duration
	sync.Mutex
}

func CreateScanConfigWatcher(scanConfigChan chan *map[string]models.ScanConfig,
	backendClient *client.ClientWithResponses,
) *ScanConfigWatcher {
	return &ScanConfigWatcher{
		existingScanConfigMap: make(map[string]models.ScanConfig),
		stopChan:              make(chan struct{}),
		scanConfigChan:        scanConfigChan,
		backendClient:         backendClient,
		Mutex:                 sync.Mutex{},
	}
}

func (scw *ScanConfigWatcher) getScanConfigs() ([]models.ScanConfig, error) {
	params := &models.GetScanConfigsParams{
		Page:     1,
		PageSize: 100,
	}

	var scanConfigs []models.ScanConfig
	paginatedScanConfigs, err := scw.getScanConfigsByPage(params)
	if err != nil {
		return scanConfigs, fmt.Errorf("failed to get scanconfigs by page:%d, pageSize: %d, error:%v", params.Page, params.PageSize, err)
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
		return &models.ScanConfigs{}, fmt.Errorf("failed to get a scan configs: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return &models.ScanConfigs{}, fmt.Errorf("no scan configs: empty body")
		}
		return resp.JSON200, nil
	default:
		message := ""
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			message = *resp.JSONDefault.Message
		}
		return &models.ScanConfigs{}, fmt.Errorf("failed to get scan configs. status code=%v: %s", resp.StatusCode(), message)
	}
}

func scanConfigSliceToMap(scanConfigSlice []models.ScanConfig) map[string]models.ScanConfig {
	scanConfigMap := make(map[string]models.ScanConfig)
	for _, scanConfig := range scanConfigSlice {
		scanConfigMap[*scanConfig.Id] = scanConfig
	}
	return scanConfigMap
}

func (scw *ScanConfigWatcher) checkNewScanConfigs() (map[string]models.ScanConfig, error) {
	scanConfigs, err := scw.getScanConfigs()
	if err != nil {
		return map[string]models.ScanConfig{}, fmt.Errorf("failed to check new scan configs: %v", err)
	}
	scanConfigMap := scanConfigSliceToMap(scanConfigs)
	newScanConfigMap := make(map[string]models.ScanConfig)
	scw.Lock()
	defer scw.Unlock()
	for k, v := range scanConfigMap {
		if _, ok := scw.existingScanConfigMap[k]; !ok {
			newScanConfigMap[k] = v
		}
	}
	scw.existingScanConfigMap = scanConfigMap
	return newScanConfigMap, nil
}

func (scw *ScanConfigWatcher) Start(errChan chan struct{}) {
	// Clear
	close(scw.stopChan)
	scw.stopChan = make(chan struct{})
	for {
		select {
		case <-time.After(scw.scanConfigWatchInterval):
			newScanConfigMap, err := scw.checkNewScanConfigs()
			if err != nil {
				if errChan != nil {
					errChan <- struct{}{}
				}
			}
			if len(newScanConfigMap) > 0 {
				scw.scanConfigChan <- &newScanConfigMap
			}
		case <-scw.stopChan:
			log.Infof("Stop watching scan configs.")
			return
		}
	}
}

func (scw *ScanConfigWatcher) Stop() {
	scw.stopChan <- struct{}{}
}
