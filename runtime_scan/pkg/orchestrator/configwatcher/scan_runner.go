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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	_scanner "github.com/openclarity/vmclarity/runtime_scan/pkg/scanner"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func (scw *ScanConfigWatcher) runNewScans(ctx context.Context, scanConfigs []models.ScanConfig) {
	for _, sc := range scanConfigs {
		scanConfig := sc
		if err := scw.scan(ctx, &scanConfig); err != nil {
			log.Errorf("Failed to schedule a scan with scan config ID=%s: %v", *scanConfig.Id, err)
		}
	}
}

func (scw *ScanConfigWatcher) scan(ctx context.Context, scanConfig *models.ScanConfig) error {
	// TODO: check if existing scan or a new scan
	targetInstances, scanID, err := scw.initNewScan(ctx, scanConfig)
	if err != nil {
		return fmt.Errorf("failed to init new scan: %v", err)
	}

	scanner := _scanner.CreateScanner(scw.scannerConfig, scw.providerClient, scw.backendClient, scanConfig, targetInstances, scanID)
	if err := scanner.Scan(ctx); err != nil {
		return fmt.Errorf("failed to scan: %v", err)
	}

	return nil
}

// initNewScan Initialized a new scan, returns target instances and scan ID.
func (scw *ScanConfigWatcher) initNewScan(ctx context.Context, scanConfig *models.ScanConfig) ([]*types.TargetInstance, string, error) {
	// Create scan in pending
	now := time.Now().UTC()
	scan := &models.Scan{
		ScanConfig: &models.ScanConfigRelationship{
			Id: *scanConfig.Id,
		},
		ScanConfigSnapshot: &models.ScanConfigData{
			Scope:              scanConfig.Scope,
			ScanFamiliesConfig: scanConfig.ScanFamiliesConfig,
		},
		StartTime: &now,
		State:     utils.PointerTo[models.ScanState](models.Pending),
		Summary: &models.ScanSummary{
			JobsCompleted:          utils.PointerTo[int](0),
			JobsLeftToRun:          utils.PointerTo[int](0),
			TotalExploits:          utils.PointerTo[int](0),
			TotalMalware:           utils.PointerTo[int](0),
			TotalMisconfigurations: utils.PointerTo[int](0),
			TotalPackages:          utils.PointerTo[int](0),
			TotalRootkits:          utils.PointerTo[int](0),
			TotalSecrets:           utils.PointerTo[int](0),
			TotalVulnerabilities: &models.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo[int](0),
				TotalHighVulnerabilities:       utils.PointerTo[int](0),
				TotalLowVulnerabilities:        utils.PointerTo[int](0),
				TotalMediumVulnerabilities:     utils.PointerTo[int](0),
				TotalNegligibleVulnerabilities: utils.PointerTo[int](0),
			},
		},
	}
	scanID, err := scw.backendClient.PostScan(ctx, *scan)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create a scan: %v", err)
	}

	// Do discovery of targets
	instances, err := scw.providerClient.Discover(ctx, scan.ScanConfigSnapshot.Scope)
	if err != nil {
		return nil, "", fmt.Errorf("failed to discover instances to scan: %v", err)
	}
	targetInstances, err := scw.createTargetInstances(ctx, instances)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create targets: %v", err)
	}

	// Move scan to discovered and add the discovered targets.
	targetIds := getTargetIDs(targetInstances)
	scan = &models.Scan{
		TargetIDs:    targetIds,
		State:        utils.PointerTo[models.ScanState](models.Discovered),
		StateMessage: utils.PointerTo[string]("Targets for scan successfully discovered"),
		// TODO sam why do we need this again?
		Summary: &models.ScanSummary{
			JobsCompleted:          utils.PointerTo[int](0),
			JobsLeftToRun:          utils.PointerTo[int](0),
			TotalExploits:          utils.PointerTo[int](0),
			TotalMalware:           utils.PointerTo[int](0),
			TotalMisconfigurations: utils.PointerTo[int](0),
			TotalPackages:          utils.PointerTo[int](0),
			TotalRootkits:          utils.PointerTo[int](0),
			TotalSecrets:           utils.PointerTo[int](0),
			TotalVulnerabilities: &models.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo[int](0),
				TotalHighVulnerabilities:       utils.PointerTo[int](0),
				TotalMediumVulnerabilities:     utils.PointerTo[int](0),
				TotalLowVulnerabilities:        utils.PointerTo[int](0),
				TotalNegligibleVulnerabilities: utils.PointerTo[int](0),
			},
		},
	}
	err = scw.backendClient.PatchScan(ctx, scanID, scan)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update scan: %v", err)
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

func (scw *ScanConfigWatcher) createTargetInstances(ctx context.Context, instances []types.Instance) ([]*types.TargetInstance, error) {
	targetInstances := make([]*types.TargetInstance, 0, len(instances))
	for i, instance := range instances {
		targetID, err := scw.createTarget(ctx, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to create target. instanceID=%v: %v", instance.GetID(), err)
		}
		targetInstances = append(targetInstances, &types.TargetInstance{
			TargetID: targetID,
			Instance: instances[i],
		})
	}

	return targetInstances, nil
}

func (scw *ScanConfigWatcher) createTarget(ctx context.Context, instance types.Instance) (string, error) {
	info := models.TargetType{}
	instanceProvider := models.AWS
	err := info.FromVMInfo(models.VMInfo{
		InstanceID:       instance.GetID(),
		InstanceProvider: &instanceProvider,
		Location:         instance.GetLocation(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create VMInfo: %v", err)
	}
	return scw.backendClient.PostTarget(ctx, models.Target{
		TargetInfo: &info,
	})
}
