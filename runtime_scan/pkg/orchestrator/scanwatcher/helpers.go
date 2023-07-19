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

package scanwatcher

import (
	"errors"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func newVulnerabilityScanSummary() *models.VulnerabilityScanSummary {
	return &models.VulnerabilityScanSummary{
		TotalCriticalVulnerabilities:   utils.PointerTo[int](0),
		TotalHighVulnerabilities:       utils.PointerTo[int](0),
		TotalMediumVulnerabilities:     utils.PointerTo[int](0),
		TotalLowVulnerabilities:        utils.PointerTo[int](0),
		TotalNegligibleVulnerabilities: utils.PointerTo[int](0),
	}
}

func newAssetScanSummary() *models.ScanFindingsSummary {
	return &models.ScanFindingsSummary{
		TotalExploits:          utils.PointerTo[int](0),
		TotalMalware:           utils.PointerTo[int](0),
		TotalMisconfigurations: utils.PointerTo[int](0),
		TotalPackages:          utils.PointerTo[int](0),
		TotalRootkits:          utils.PointerTo[int](0),
		TotalSecrets:           utils.PointerTo[int](0),
		TotalVulnerabilities:   newVulnerabilityScanSummary(),
	}
}

func newAssetScanFromScan(scan *models.Scan, assetID string) (*models.AssetScan, error) {
	if scan == nil {
		return nil, errors.New("failed to create AssetScan: Scan is nil")
	}

	if scan.AssetScanTemplate == nil || scan.AssetScanTemplate.ScanFamiliesConfig == nil {
		return nil, errors.New("failed to create AssetScan: AssetScanTemplate and/or AssetScanTemplate.ScanFamiliesConfig is nil")
	}
	familiesConfig := scan.AssetScanTemplate.ScanFamiliesConfig

	return &models.AssetScan{
		Summary: newAssetScanSummary(),
		Scan: &models.ScanRelationship{
			Id: *scan.Id,
		},
		Asset: &models.AssetRelationship{
			Id: assetID,
		},
		ScanFamiliesConfig:            familiesConfig,
		ScannerInstanceCreationConfig: scan.AssetScanTemplate.ScannerInstanceCreationConfig,
		Status: &models.AssetScanStatus{
			Exploits: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Exploits),
			},
			General: &models.AssetScanState{
				Errors: nil,
				State:  utils.PointerTo(models.AssetScanStateStatePending),
			},
			Malware: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Malware),
			},
			Misconfigurations: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Misconfigurations),
			},
			Rootkits: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Rootkits),
			},
			Sbom: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Sbom),
			},
			Secrets: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Secrets),
			},
			Vulnerabilities: &models.AssetScanState{
				Errors: nil,
				State:  getInitStateFromFamilyConfig(familiesConfig.Vulnerabilities),
			},
		},
		ResourceCleanup: utils.PointerTo(models.ResourceCleanupStatePending),
	}, nil
}

func getInitStateFromFamilyConfig(config models.FamilyConfigEnabler) *models.AssetScanStateState {
	if config == nil || !config.IsEnabled() {
		return utils.PointerTo(models.AssetScanStateStateNotScanned)
	}

	return utils.PointerTo(models.AssetScanStateStatePending)
}

func newScanSummary() *models.ScanSummary {
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

func updateScanSummaryFromAssetScan(scan *models.Scan, result models.AssetScan) error {
	if result.Summary == nil {
		return errors.New("invalid AssetScan: Summary field is nil")
	}

	if scan.Summary == nil {
		scan.Summary = newScanSummary()
	}

	state, ok := result.GetGeneralState()
	if !ok {
		return fmt.Errorf("general state must not be nil for AssetScan. AssetScanID=%s", *result.Id)
	}

	s, r := scan.Summary, result.Summary

	switch state {
	case models.AssetScanStateStateNotScanned:
	case models.AssetScanStateStatePending, models.AssetScanStateStateScheduled, models.AssetScanStateStateReadyToScan:
		fallthrough
	case models.AssetScanStateStateInProgress, models.AssetScanStateStateAborted:
		s.JobsLeftToRun = utils.PointerTo(*s.JobsLeftToRun + 1)
	case models.AssetScanStateStateDone:
		s.JobsCompleted = utils.PointerTo(*s.JobsCompleted + 1)
		s.TotalExploits = utils.PointerTo(*s.TotalExploits + *r.TotalExploits)
		s.TotalMalware = utils.PointerTo(*s.TotalMalware + *r.TotalMalware)
		s.TotalMisconfigurations = utils.PointerTo(*s.TotalMisconfigurations + *r.TotalMisconfigurations)
		s.TotalPackages = utils.PointerTo(*s.TotalPackages + *r.TotalPackages)
		s.TotalRootkits = utils.PointerTo(*s.TotalRootkits + *r.TotalRootkits)
		s.TotalSecrets = utils.PointerTo(*s.TotalSecrets + *r.TotalSecrets)
		s.TotalVulnerabilities = &models.VulnerabilityScanSummary{
			TotalCriticalVulnerabilities: utils.PointerTo(*s.TotalVulnerabilities.TotalCriticalVulnerabilities +
				*r.TotalVulnerabilities.TotalCriticalVulnerabilities),
			TotalHighVulnerabilities: utils.PointerTo(*s.TotalVulnerabilities.TotalHighVulnerabilities +
				*r.TotalVulnerabilities.TotalHighVulnerabilities),
			TotalLowVulnerabilities: utils.PointerTo(*s.TotalVulnerabilities.TotalLowVulnerabilities +
				*r.TotalVulnerabilities.TotalLowVulnerabilities),
			TotalMediumVulnerabilities: utils.PointerTo(*s.TotalVulnerabilities.TotalMediumVulnerabilities +
				*r.TotalVulnerabilities.TotalMediumVulnerabilities),
			TotalNegligibleVulnerabilities: utils.PointerTo(*s.TotalVulnerabilities.TotalCriticalVulnerabilities +
				*r.TotalVulnerabilities.TotalNegligibleVulnerabilities),
		}
	}

	return nil
}
