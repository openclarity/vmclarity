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

package scheduler

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
	_scanner "github.com/openclarity/vmclarity/runtime_scan/pkg/scanner"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

type Scheduler struct {
	stopChan        chan struct{}
	scanConfigsChan chan []models.ScanConfig
	scannerConfig   *_config.ScannerConfig
	providerClient  provider.Client
	backendClient   *client.ClientWithResponses
}

type Params struct {
	Interval   time.Duration
	StartTime  time.Time
	SingleScan bool
}

func CreateScheduler(scanConfigsChan chan []models.ScanConfig,
	scannerConfig *_config.ScannerConfig,
	providerClient provider.Client,
	backendClient *client.ClientWithResponses,
) *Scheduler {
	return &Scheduler{
		stopChan:        make(chan struct{}),
		scanConfigsChan: scanConfigsChan,
		scannerConfig:   scannerConfig,
		providerClient:  providerClient,
		backendClient:   backendClient,
	}
}

func (s *Scheduler) Start() {
	// Clear
	close(s.stopChan)
	s.stopChan = make(chan struct{})
	go func() {
		for {
			select {
			case scanConfigs := <-s.scanConfigsChan:
				s.scheduleNewScans(scanConfigs)
			case <-s.stopChan:
				log.Infof("Stop scheduling scans.")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) scheduleNewScans(scanConfigs []models.ScanConfig) {
	for _, scanConfig := range scanConfigs {
		scanConfig := scanConfig
		// Now only SingleScheduledScanConfigs will be started, so don't need to real schedule.
		if err := s.scan(context.Background(), &scanConfig); err != nil {
			log.Errorf("falied to schedule a scan with scan config ID=%s: %v", *scanConfig.Id, err)
		}
	}
}

func (s *Scheduler) scan(ctx context.Context, scanConfig *models.ScanConfig) error {
	// TODO: check if existing scan or a new scan
	targetInstances, scanID, err := s.initNewScan(ctx, scanConfig)
	if err != nil {
		return fmt.Errorf("failed to init new scan: %v", err)
	}

	scanner := _scanner.CreateScanner(s.scannerConfig, s.providerClient, s.backendClient, scanConfig, targetInstances, scanID)
	scanDone := make(chan struct{})
	if err := scanner.Scan(ctx, scanDone); err != nil {
		return fmt.Errorf("failed to scan: %v", err)
	}

	return nil
}

// initNewScan Initialized a new scan, returns target instances and scan ID.
func (s *Scheduler) initNewScan(ctx context.Context, scanConfig *models.ScanConfig) ([]*types.TargetInstance, string, error) {
	instances, err := s.providerClient.Discover(ctx, scanConfig.Scope)
	if err != nil {
		return nil, "", fmt.Errorf("failed to discover instances to scan: %v", err)
	}

	targetInstances, err := s.createTargetInstances(ctx, instances)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create targets: %v", err)
	}

	now := time.Now().UTC()
	scan := &models.Scan{
		ScanConfigId:       scanConfig.Id,
		ScanFamiliesConfig: scanConfig.ScanFamiliesConfig,
		StartTime:          &now,
		TargetIDs:          getTargetIDs(targetInstances),
	}
	scanID, err := s.getOrCreateScan(ctx, scan)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create a scan: %v", err)
	}

	return targetInstances, scanID, nil
}

func getTargetIDs(targetInstances []*types.TargetInstance) *[]string {
	ret := make([]string, len(targetInstances))
	for i, targetInstance := range targetInstances {
		ret[i] = targetInstance.TargetID
	}

	return &ret
}

func (s *Scheduler) createTargetInstances(ctx context.Context, instances []types.Instance) ([]*types.TargetInstance, error) {
	targetInstances := make([]*types.TargetInstance, 0, len(instances))
	for i, instance := range instances {
		target, err := s.getOrCreateTarget(ctx, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to get or create target. instanceID=%v: %v", instance.GetID(), err)
		}
		targetInstances = append(targetInstances, &types.TargetInstance{
			TargetID: *target.Id,
			Instance: instances[i],
		})
	}

	return targetInstances, nil
}

func (s *Scheduler) getOrCreateTarget(ctx context.Context, instance types.Instance) (*models.Target, error) {
	info := models.TargetType{}
	instanceProvider := models.AWS
	err := info.FromVMInfo(models.VMInfo{
		InstanceID:       utils.StringPtr(instance.GetID()),
		InstanceProvider: &instanceProvider,
		Location:         utils.StringPtr(instance.GetLocation()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create VMInfo: %v", err)
	}
	resp, err := s.backendClient.PostTargetsWithResponse(ctx, models.Target{
		TargetInfo: &info,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to post target: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return nil, fmt.Errorf("failed to create a target: empty body")
		}
		return resp.JSON201, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return nil, fmt.Errorf("failed to create a target: empty body on conflict")
		}
		return resp.JSON409, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to post target. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to post target. status code=%v", resp.StatusCode())
	}
}

// nolint:cyclop
func (s *Scheduler) getOrCreateScan(ctx context.Context, scan *models.Scan) (string, error) {
	resp, err := s.backendClient.PostScansWithResponse(ctx, *scan)
	if err != nil {
		return "", fmt.Errorf("failed to post a scan: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body")
		}
		if resp.JSON201.Id == nil {
			return "", fmt.Errorf("scan id is nil")
		}
		return *resp.JSON201.Id, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body on conflict")
		}
		if resp.JSON409.Id == nil {
			return "", fmt.Errorf("scan id on conflict is nil")
		}
		return *resp.JSON409.Id, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return "", fmt.Errorf("failed to post scan. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return "", fmt.Errorf("failed to post scan. status code=%v", resp.StatusCode())
	}
}
