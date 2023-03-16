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

	"github.com/aptible/supercronic/cronexpr"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

const (
	timeWindow = 5 * time.Minute
)

type ScanConfigWatcher struct {
	backendClient  *backendclient.BackendClient
	providerClient provider.Client
	scannerConfig  *_config.ScannerConfig
}

func CreateScanConfigWatcher(
	backendClient *backendclient.BackendClient,
	providerClient provider.Client,
	scannerConfig _config.ScannerConfig,
) *ScanConfigWatcher {
	return &ScanConfigWatcher{
		backendClient:  backendClient,
		providerClient: providerClient,
		scannerConfig:  &scannerConfig,
	}
}

func (scw *ScanConfigWatcher) hasRunningScan(scanConfigID string) (bool, error) {
	// We want to check if there is an existing scan related to the given scan config ID that is still running (not done or failed).
	odataFilter := fmt.Sprintf("scanConfig/id eq '%s' and state ne '%s' and state ne '%s'",
		scanConfigID, models.ScanStateDone, models.ScanStateFailed)
	params := models.GetScansParams{
		Filter: &odataFilter,
		Count:  utils.PointerTo(true),
	}

	scans, err := scw.backendClient.GetScans(context.TODO(), params)
	if err != nil {
		return false, fmt.Errorf("failed to get a scans: %v", err)
	}

	return *scans.Count > 0, nil
}

func (scw *ScanConfigWatcher) getScanConfigsToScan() ([]models.ScanConfig, error) {
	scanConfigsToScan := make([]models.ScanConfig, 0)
	// Get all enabled scan configs
	scanConfigs, err := scw.backendClient.GetScanConfigs(context.TODO(), models.GetScanConfigsParams{
		Filter: utils.PointerTo(fmt.Sprintf("disabled eq null or disabled eq false")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check new scan configs: %v", err)
	}

	log.Infof("Got %d enabled ScanConfigs objects.", len(*scanConfigs.Items))

	now := time.Now()
	for _, scanConfig := range *scanConfigs.Items {
		shouldScan := false
		scanConfigID := *scanConfig.Id
		operationTime := *scanConfig.Scheduled.OperationTime
		shouldScan, err = scw.shouldScan(scanConfigID, operationTime, now)
		if err != nil {
			log.Errorf("Failed to get scans for scan config ID=%s: %v", scanConfigID, err)
			continue
		}

		if shouldScan {
			log.Infof("A new scan should be started from ScanConfig %s", scanConfigID)
			scanConfigsToScan = append(scanConfigsToScan, scanConfig)
			if scanConfig.Scheduled.CronLine != nil {
				// calculate next operation time based on current operation time
				nextOperationTime := cronexpr.MustParse(*scanConfig.Scheduled.CronLine).Next(operationTime)
				scanConfig.Scheduled.OperationTime = &nextOperationTime
				log.Debugf("Patching ScanConfig %s with a new operation time (%s)", scanConfigID, nextOperationTime.String())
			} else {
				// not a periodic scan, we should disable the scan config, so it will not be fetched again.
				scanConfig.Disabled = utils.PointerTo(false)
				log.Debugf("Patching ScanConfig %s with disabled (%v)", scanConfigID, *scanConfig.Disabled)
			}
			if err = scw.backendClient.PatchScanConfig(context.TODO(), scanConfigID, &scanConfig); err != nil {
				log.Errorf("Failed to patch scan config: %v", err)
			}
		} else {
			log.Debugf("No scan should be started from ScanConfig %s", scanConfigID)
			if operationTime.Before(now) && scanConfig.Scheduled.CronLine != nil {
				// If operationTime is not within the window, and it was in the past,
				// we will calculate the next operation time until we will find one that is in the future.
				nextOperationTime := findFirstOperationTimeInTheFuture(operationTime, now, *scanConfig.Scheduled.CronLine)
				scanConfig.Scheduled.OperationTime = &nextOperationTime
				log.Debugf("Patching ScanConfig %s with a new operation time (%s)", scanConfigID, nextOperationTime.String())
				if err = scw.backendClient.PatchScanConfig(context.TODO(), scanConfigID, &scanConfig); err != nil {
					log.Errorf("Failed to patch scan config: %v", err)
				}
			}
		}
	}

	return scanConfigsToScan, nil
}

func findFirstOperationTimeInTheFuture(operationTime time.Time, now time.Time, cronLine string) time.Time {
	expr := cronexpr.MustParse(cronLine)
	for operationTime.Before(now) {
		operationTime = expr.Next(operationTime)
	}
	return operationTime
}

// isWithinTheWindow checks if `checkTime` is within the window (after `now` and before `now + window`).
func isWithinTheWindow(checkTime, now time.Time, window time.Duration) bool {
	if checkTime.Before(now) {
		return false
	}

	endWindowTime := now.Add(window)
	return checkTime.Before(endWindowTime)
}

func (scw *ScanConfigWatcher) shouldScan(scanConfigID string, operationTime time.Time, now time.Time) (bool, error) {
	// Skip processing ScanConfig because its operationTime is not within the start window
	if !isWithinTheWindow(operationTime, now, timeWindow) {
		log.Debugf("ScanConfig %s start time %v outside of the start window %v - %v",
			scanConfigID, operationTime.Format(time.RFC3339), now.Format(time.RFC3339), now.Add(timeWindow).Format(time.RFC3339))
		return false, nil
	}

	// Check running scans for specific scan config
	hasRunningScan, err := scw.hasRunningScan(scanConfigID)
	if err != nil {
		return false, fmt.Errorf("failed to check if there are running or completed scans: %v", err)
	}
	if hasRunningScan {
		log.Debugf("ScanConfig %s has a running scan", scanConfigID)
	}

	// If operation time is withing the window and there is no running scan we should run a scan.
	return !hasRunningScan, nil
}

func (scw *ScanConfigWatcher) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(scw.scannerConfig.ScanConfigWatchInterval):
				log.Debug("Looking for ScanConfigs to Scan")
				// nolint:contextcheck
				scanConfigsToScan, err := scw.getScanConfigsToScan()
				if err != nil {
					log.Warnf("Failed to get scan configs to scan: %v", err)
					break
				}
				log.Debugf("Found %d ScanConfigs that need to start scanning.", len(scanConfigsToScan))
				if len(scanConfigsToScan) > 0 {
					scw.runNewScans(ctx, scanConfigsToScan)
				}
			case <-ctx.Done():
				log.Infof("Stop watching scan configs.")
				return
			}
		}
	}()
}
