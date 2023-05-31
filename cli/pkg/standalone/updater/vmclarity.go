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

package updater

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type VMClarityUpdater struct {
	client       *backendclient.BackendClient
	scanID       string
	scanResultID string
}

func NewVMClarityUpdater(client *backendclient.BackendClient, scanID, scanResultID string) (*VMClarityUpdater, error) {
	if client == nil {
		return nil, errors.New("backend client must not be nil")
	}

	return &VMClarityUpdater{
		client:       client,
		scanID:       scanID,
		scanResultID: scanResultID,
	}, nil
}

func (u *VMClarityUpdater) SetScanIDIfNeeded(ctx context.Context) error {
	if u.scanID != "" {
		return nil
	}
	scanResult, err := u.client.GetScanResult(ctx, u.scanResultID, models.GetScanResultsScanResultIDParams{})
	if err != nil {
		return fmt.Errorf("failed to get scan result by ID=%s: %v", u.scanResultID, err)
	}
	u.scanID = scanResult.Scan.Id

	return nil
}

func (u *VMClarityUpdater) UpdateScanStateAndSummary(ctx context.Context) error {
	scan, err := u.updatedScanSummary(ctx)
	if err != nil {
		return fmt.Errorf("failed to update scan results summary: %v", err)
	}

	scan.EndTime = utils.PointerTo(time.Now())
	scan.State = utils.PointerTo(models.ScanStateDone)
	scan.StateMessage = utils.PointerTo(utils.AllScanJobsCompletedMessage)
	scan.StateReason = utils.PointerTo(models.ScanStateReasonSuccess)

	err = u.client.PatchScan(ctx, u.scanID, scan)
	if err != nil {
		return fmt.Errorf("failed to patch the scan ID=%s: %v", u.scanID, err)
	}

	return nil
}

func (u *VMClarityUpdater) updatedScanSummary(ctx context.Context) (*models.Scan, error) {
	scan, err := u.client.GetScan(ctx, u.scanID, models.GetScansScanIDParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get scan to update status: %v", err)
	}

	scanResultSummary, err := u.client.GetScanResultSummary(ctx, u.scanResultID)
	if err != nil {
		return nil, fmt.Errorf("failed to get result summary to update status: %v", err)
	}

	// Update the scan summary with the summary from the completed scan result
	scan.Summary.JobsCompleted = utils.PointerTo(*scan.Summary.JobsCompleted + 1)
	scan.Summary.JobsLeftToRun = utils.PointerTo(*scan.Summary.JobsLeftToRun - 1)
	scan.Summary.TotalExploits = utils.PointerTo(*scan.Summary.TotalExploits + *scanResultSummary.TotalExploits)
	scan.Summary.TotalMalware = utils.PointerTo(*scan.Summary.TotalMalware + *scanResultSummary.TotalMalware)
	scan.Summary.TotalMisconfigurations = utils.PointerTo(*scan.Summary.TotalMisconfigurations + *scanResultSummary.TotalMisconfigurations)
	scan.Summary.TotalPackages = utils.PointerTo(*scan.Summary.TotalPackages + *scanResultSummary.TotalPackages)
	scan.Summary.TotalRootkits = utils.PointerTo(*scan.Summary.TotalRootkits + *scanResultSummary.TotalRootkits)
	scan.Summary.TotalSecrets = utils.PointerTo(*scan.Summary.TotalSecrets + *scanResultSummary.TotalSecrets)
	scan.Summary.TotalVulnerabilities = &models.VulnerabilityScanSummary{
		TotalCriticalVulnerabilities:   utils.PointerTo(*scan.Summary.TotalVulnerabilities.TotalCriticalVulnerabilities + *scanResultSummary.TotalVulnerabilities.TotalCriticalVulnerabilities),
		TotalHighVulnerabilities:       utils.PointerTo(*scan.Summary.TotalVulnerabilities.TotalHighVulnerabilities + *scanResultSummary.TotalVulnerabilities.TotalHighVulnerabilities),
		TotalLowVulnerabilities:        utils.PointerTo(*scan.Summary.TotalVulnerabilities.TotalLowVulnerabilities + *scanResultSummary.TotalVulnerabilities.TotalLowVulnerabilities),
		TotalMediumVulnerabilities:     utils.PointerTo(*scan.Summary.TotalVulnerabilities.TotalMediumVulnerabilities + *scanResultSummary.TotalVulnerabilities.TotalMediumVulnerabilities),
		TotalNegligibleVulnerabilities: utils.PointerTo(*scan.Summary.TotalVulnerabilities.TotalCriticalVulnerabilities + *scanResultSummary.TotalVulnerabilities.TotalNegligibleVulnerabilities),
	}

	return scan, nil
}
