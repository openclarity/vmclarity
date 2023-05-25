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
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	_scanner "github.com/openclarity/vmclarity/runtime_scan/pkg/scanner"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

func (scw *ScanConfigWatcher) scan(ctx context.Context, scanConfig *models.ScanConfig) error {
	// TODO: check if existing scan or a new scan
	targets, scanID, err := scw.initNewScan(ctx, scanConfig)
	if err != nil {
		return fmt.Errorf("failed to init new scan: %v", err)
	}

	scanner := _scanner.CreateScanner(scw.scannerConfig, scw.providerClient, scw.backendClient, scanConfig, targets, scanID)
	go scanner.Scan(ctx)

	return nil
}

// initNewScan Initialized a new scan, returns target instances and scan ID.
func (scw *ScanConfigWatcher) initNewScan(ctx context.Context, scanConfig *models.ScanConfig) ([]models.Target, string, error) {
	// Create scan in pending
	now := time.Now().UTC()
	scan := &models.Scan{
		ScanConfig: &models.ScanConfigRelationship{
			Id: *scanConfig.Id,
		},
		ScanConfigSnapshot: &models.ScanConfigSnapshot{
			MaxParallelScanners: scanConfig.MaxParallelScanners,
			Name:                scanConfig.Name,
			ScanFamiliesConfig:  scanConfig.ScanFamiliesConfig,
			Scheduled:           scanConfig.Scheduled,
			Scope:               scanConfig.Scope,
		},
		StartTime: &now,
		State:     utils.PointerTo(models.ScanStatePending),
		Summary:   createInitScanSummary(),
	}
	var scanID string
	createdScan, err := scw.backendClient.PostScan(ctx, *scan)
	if err != nil {
		var conErr backendclient.ScanConflictError
		if errors.As(err, &conErr) {
			log.Infof("Scan already exist. scan id=%v.", *conErr.ConflictingScan.Id)
			scanID = *conErr.ConflictingScan.Id
		} else {
			return nil, "", fmt.Errorf("failed to post scan: %v", err)
		}
	} else {
		scanID = *createdScan.Id
	}

	// Do discovery of targets
	targets, err := scw.backendClient.GetTargets(ctx, models.GetTargetsParams{
		Filter: createdScan.ScanConfigSnapshot.Scope,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to query targets to scan: %w", err)
	}

	// We just want the IDs right now
	targetIds := []string{}
	for _, target := range *targets.Items {
		targetIds = append(targetIds, *target.Id)
	}

	// Move scan to discovered and add the discovered targets.
	scan = &models.Scan{
		TargetIDs:    &targetIds,
		State:        utils.PointerTo(models.ScanStateDiscovered),
		StateMessage: utils.PointerTo("Targets for scan successfully discovered"),
	}
	err = scw.backendClient.PatchScan(ctx, scanID, scan)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update scan: %v", err)
	}

	return *targets.Items, scanID, nil
}

func createInitScanSummary() *models.ScanSummary {
	return &models.ScanSummary{
		JobsCompleted:          utils.PointerTo(0),
		JobsLeftToRun:          utils.PointerTo(0),
		TotalExploits:          utils.PointerTo(0),
		TotalMalware:           utils.PointerTo(0),
		TotalMisconfigurations: utils.PointerTo(0),
		TotalPackages:          utils.PointerTo(0),
		TotalRootkits:          utils.PointerTo(0),
		TotalSecrets:           utils.PointerTo(0),
		TotalVulnerabilities: &models.VulnerabilityScanSummary{
			TotalCriticalVulnerabilities:   utils.PointerTo(0),
			TotalHighVulnerabilities:       utils.PointerTo(0),
			TotalLowVulnerabilities:        utils.PointerTo(0),
			TotalMediumVulnerabilities:     utils.PointerTo(0),
			TotalNegligibleVulnerabilities: utils.PointerTo(0),
		},
	}
}
