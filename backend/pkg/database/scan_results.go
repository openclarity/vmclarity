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

	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
)

// TODO after db design.
type ScanResults struct {
	ID               string
	Sbom             *SbomScanResults
	Vulnerability    *VulnerabilityScanResults
	Malware          *MalwareScanResults
	Rootkit          *RootkitScanScanResults
	Secret           *SecretScanResults
	Misconfiguration *MisconfigurationScanResults
	Exploit          *ExploitScanResults
}

// TODO after db design.
type SbomScanResults struct {
	Results models.SbomScan
}

type VulnerabilityScanResults struct {
	Results models.VulnerabilityScan
}

type MalwareScanResults struct {
	Results models.MalwareScan
}

type RootkitScanScanResults struct {
	Results models.RootkitScan
}

type SecretScanResults struct {
	Results models.SecretScan
}

type MisconfigurationScanResults struct {
	Results models.MisconfigurationScan
}

type ExploitScanResults struct {
	Results models.ExploitScan
}

//nolint:interfacebloat
type ScanResultsTable interface {
	ListScanResults(targetID models.TargetID, params models.GetTargetsTargetIDScanResultsParams) (*[]models.ScanResults, error)
	CreateScanResults(targetID models.TargetID, scanResults *ScanResults) (*models.ScanResultsSummary, error)
	GetScanResults(targetID models.TargetID, scanID models.ScanID) (*models.ScanResults, error)
	GetSBOM(targetID models.TargetID, scanID models.ScanID) (*models.SbomScan, error)
	GetVulnerabilities(targetID models.TargetID, scanID models.ScanID) (*models.VulnerabilityScan, error)
	GetMalwares(targetID models.TargetID, scanID models.ScanID) (*models.MalwareScan, error)
	GetRootkits(targetID models.TargetID, scanID models.ScanID) (*models.RootkitScan, error)
	GetSecrets(targetID models.TargetID, scanID models.ScanID) (*models.SecretScan, error)
	GetMisconfigurations(targetID models.TargetID, scanID models.ScanID) (*models.MisconfigurationScan, error)
	GetExploits(targetID models.TargetID, scanID models.ScanID) (*models.ExploitScan, error)
	UpdateScanResults(targetID models.TargetID, scanID models.ScanID, scanResults *ScanResults) (*models.ScanResultsSummary, error)
}

type ScanResultsTableHandler struct {
	db *gorm.DB
}

func (db *Handler) ScanResultsTable() ScanResultsTable {
	return &ScanResultsTableHandler{
		db: db.DB,
	}
}

func (s *ScanResultsTableHandler) ListScanResults(targetID models.TargetID, params models.GetTargetsTargetIDScanResultsParams,
) (*[]models.ScanResults, error) {
	return &[]models.ScanResults{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) CreateScanResults(targetID models.TargetID, scanResults *ScanResults,
) (*models.ScanResultsSummary, error) {
	return &models.ScanResultsSummary{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetScanResults(targetID models.TargetID, scanID models.ScanID) (*models.ScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetSBOM(targetID models.TargetID, scanID models.ScanID) (*models.SbomScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetVulnerabilities(targetID models.TargetID, scanID models.ScanID) (*models.VulnerabilityScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetMalwares(targetID models.TargetID, scanID models.ScanID) (*models.MalwareScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetRootkits(targetID models.TargetID, scanID models.ScanID) (*models.RootkitScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetSecrets(targetID models.TargetID, scanID models.ScanID) (*models.SecretScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetMisconfigurations(targetID models.TargetID, scanID models.ScanID) (*models.MisconfigurationScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetExploits(targetID models.TargetID, scanID models.ScanID) (*models.ExploitScan, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) UpdateScanResults(
	targetID models.TargetID,
	scanID models.ScanID,
	scanResults *ScanResults,
) (*models.ScanResultsSummary, error) {
	return &models.ScanResultsSummary{}, fmt.Errorf("not implemented")
}

// TODO after db design.
func CreateDBScanResultsFromModel(scanResults *models.ScanResults) *ScanResults {
	return &ScanResults{
		ID: *scanResults.Id,
		Sbom: &SbomScanResults{
			Results: *scanResults.Sboms,
		},
		Vulnerability: &VulnerabilityScanResults{
			Results: *scanResults.Vulnerabilities,
		},
		Malware: &MalwareScanResults{
			Results: *scanResults.Malwares,
		},
		Rootkit: &RootkitScanScanResults{
			Results: *scanResults.Rootkits,
		},
		Secret: &SecretScanResults{
			Results: *scanResults.Secrets,
		},
		Misconfiguration: &MisconfigurationScanResults{
			Results: *scanResults.Misconfigurations,
		},
		Exploit: &ExploitScanResults{
			Results: *scanResults.Exploits,
		},
	}
}

func CreateModelScanResultsFromDB(scanResults *ScanResults) *models.ScanResults {
	return &models.ScanResults{
		Id:                &scanResults.ID,
		Sboms:             &scanResults.Sbom.Results,
		Vulnerabilities:   &scanResults.Vulnerability.Results,
		Malwares:          &scanResults.Malware.Results,
		Rootkits:          &scanResults.Rootkit.Results,
		Secrets:           &scanResults.Secret.Results,
		Misconfigurations: &scanResults.Misconfiguration.Results,
		Exploits:          &scanResults.Exploit.Results,
	}
}

func CreateModelScanResultsSummaryFromDB(scanResults *ScanResults) *models.ScanResultsSummary {
	packagesCount := len(*scanResults.Sbom.Results.Packages)
	vulnerabilitiesCount := len(*scanResults.Vulnerability.Results.Vulnerabilities)
	malwareCount := len(*scanResults.Malware.Results.Malwares)
	secretCount := len(*scanResults.Secret.Results.Secrets)
	rootkitCount := len(*scanResults.Rootkit.Results.Rootkits)
	misconfigurationCount := len(*scanResults.Misconfiguration.Results.Misconfigurations)
	exploitsCount := len(*scanResults.Exploit.Results.Exploits)
	return &models.ScanResultsSummary{
		PackagesCount:          &packagesCount,
		VulnerabilitiesCount:   &vulnerabilitiesCount,
		MalwaresCount:          &malwareCount,
		SecretsCount:           &secretCount,
		RootkitsCount:          &rootkitCount,
		MisconfigurationsCount: &misconfigurationCount,
		ExploitsCount:          &exploitsCount,
	}
}

func CountScanResultsSummary(prevSummary, newSummary *models.ScanResultsSummary) *models.ScanResultsSummary {
	packagesCount := *prevSummary.PackagesCount + *newSummary.PackagesCount
	vulnerabilitiesCount := *prevSummary.VulnerabilitiesCount + *newSummary.VulnerabilitiesCount
	malwareCount := *prevSummary.MalwaresCount + *newSummary.MalwaresCount
	secretCount := *prevSummary.SecretsCount + *newSummary.SecretsCount
	rootkitCount := *prevSummary.RootkitsCount + *newSummary.RootkitsCount
	misconfigurationCount := *prevSummary.MisconfigurationsCount + *newSummary.MisconfigurationsCount
	exploitsCount := *prevSummary.ExploitsCount + *newSummary.ExploitsCount
	return &models.ScanResultsSummary{
		PackagesCount:          &packagesCount,
		VulnerabilitiesCount:   &vulnerabilitiesCount,
		MalwaresCount:          &malwareCount,
		SecretsCount:           &secretCount,
		RootkitsCount:          &rootkitCount,
		MisconfigurationsCount: &misconfigurationCount,
		ExploitsCount:          &exploitsCount,
	}
}
