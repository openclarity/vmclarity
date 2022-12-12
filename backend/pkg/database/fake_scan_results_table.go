// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package database

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

func (fs *FakeScanResultsTable) ListScanResults(targetID models.TargetID, params models.GetTargetsTargetIDScanResultsParams,
) (*[]models.ScanResults, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	scanResults := make([]models.ScanResults, 0)
	results := *fs.scanResults
	for _, res := range results {
		scanResult := CreateModelScanResultsFromDB(res)
		scanResults = append(scanResults, *scanResult)
	}
	return &scanResults, nil
}

func (fs *FakeScanResultsTable) CreateScanResults(targetID models.TargetID, scanResults *ScanResults,
) (*models.ScanResultsSummary, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	targets[targetID].ScanResults = append(targets[targetID].ScanResults, scanResults.ID)

	scanRes := *fs.scanResults
	scanRes[scanResults.ID] = scanResults
	fs.scanResults = &scanRes
	return CreateModelScanResultsSummaryFromDB(scanResults), nil
}

func (fs *FakeScanResultsTable) GetScanResultsSummary(targetID models.TargetID, scanID models.ScanID) (*models.ScanResultsSummary, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return CreateModelScanResultsSummaryFromDB(results[scanID]), nil
}

func (fs *FakeScanResultsTable) GetSBOM(targetID models.TargetID, scanID models.ScanID) (*models.SbomScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Sbom.Results, nil
}

func (fs *FakeScanResultsTable) GetVulnerabilities(targetID models.TargetID, scanID models.ScanID) (*models.VulnerabilityScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Vulnerability.Results, nil
}

func (fs *FakeScanResultsTable) GetMalwares(targetID models.TargetID, scanID models.ScanID) (*models.MalwareScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Malware.Results, nil
}

func (fs *FakeScanResultsTable) GetRootkits(targetID models.TargetID, scanID models.ScanID) (*models.RootkitScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Rootkit.Results, nil
}

func (fs *FakeScanResultsTable) GetSecrets(targetID models.TargetID, scanID models.ScanID) (*models.SecretScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Secret.Results, nil
}

func (fs *FakeScanResultsTable) GetMisconfigurations(targetID models.TargetID, scanID models.ScanID) (*models.MisconfigurationScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Misconfiguration.Results, nil
}

func (fs *FakeScanResultsTable) GetExploits(targetID models.TargetID, scanID models.ScanID) (*models.ExploitScan, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	return &results[scanID].Exploit.Results, nil
}

func (fs *FakeScanResultsTable) UpdateScanResults(
	targetID models.TargetID,
	scanID models.ScanID,
	scanResults *ScanResults,
) (*models.ScanResultsSummary, error) {
	targets := *fs.targets
	if _, ok := targets[targetID]; !ok {
		return nil, fmt.Errorf("target not exists with ID: %s", targetID)
	}
	if !contains(scanID, targets[targetID].ScanResults) {
		return nil, fmt.Errorf("scanID %s not exists for target with ID: %s", scanID, targetID)
	}
	results := *fs.scanResults
	results[scanID] = scanResults
	fs.scanResults = &results
	return CreateModelScanResultsSummaryFromDB(results[scanID]), nil
}

func contains(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
