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

	"github.com/google/uuid"
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

type ScanResultsSummary struct {
	PackagesCount          int
	VulnerabilitiesCount   int
	MalwaresCount          int
	SecretsCount           int
	RootkitsCount          int
	MisconfigurationsCount int
	ExploitsCount          int
}

//nolint:interfacebloat
type ScanResultsTable interface {
	ListScanResults(targetID models.TargetID, params models.GetTargetsTargetIDScanResultsParams) ([]ScanResults, error)
	CreateScanResults(targetID models.TargetID, scanResults *ScanResults) (*ScanResultsSummary, error)
	GetScanResults(targetID models.TargetID, scanID models.ScanID) (*ScanResults, error)
	GetSBOM(targetID models.TargetID, scanID models.ScanID) (*SbomScanResults, error)
	GetVulnerabilities(targetID models.TargetID, scanID models.ScanID) (*VulnerabilityScanResults, error)
	GetMalwares(targetID models.TargetID, scanID models.ScanID) (*MalwareScanResults, error)
	GetRootkits(targetID models.TargetID, scanID models.ScanID) (*RootkitScanScanResults, error)
	GetSecrets(targetID models.TargetID, scanID models.ScanID) (*SecretScanResults, error)
	GetMisconfigurations(targetID models.TargetID, scanID models.ScanID) (*MisconfigurationScanResults, error)
	GetExploits(targetID models.TargetID, scanID models.ScanID) (*ExploitScanResults, error)
	UpdateScanResults(targetID models.TargetID, scanID models.ScanID, scanResults *ScanResults) (*ScanResultsSummary, error)
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
) ([]ScanResults, error) {
	return []ScanResults{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) CreateScanResults(targetID models.TargetID, scanResults *ScanResults,
) (*ScanResultsSummary, error) {
	return &ScanResultsSummary{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetScanResults(targetID models.TargetID, scanID models.ScanID) (*ScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetSBOM(targetID models.TargetID, scanID models.ScanID) (*SbomScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetVulnerabilities(targetID models.TargetID, scanID models.ScanID) (*VulnerabilityScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetMalwares(targetID models.TargetID, scanID models.ScanID) (*MalwareScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetRootkits(targetID models.TargetID, scanID models.ScanID) (*RootkitScanScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetSecrets(targetID models.TargetID, scanID models.ScanID) (*SecretScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetMisconfigurations(targetID models.TargetID, scanID models.ScanID) (*MisconfigurationScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) GetExploits(targetID models.TargetID, scanID models.ScanID) (*ExploitScanResults, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) UpdateScanResults(
	targetID models.TargetID,
	scanID models.ScanID,
	scanResults *ScanResults,
) (*ScanResultsSummary, error) {
	return &ScanResultsSummary{}, fmt.Errorf("not implemented")
}

// TODO after db design.
func CreateDBScanResultsFromModel(scanResults *models.ScanResults) *ScanResults {
	var scanResultID string
	if scanResults.Id == nil || *scanResults.Id == "" {
		scanResultID = generateScanResultsID()
	} else {
		scanResultID = *scanResults.Id
	}
	var sbomRes *SbomScanResults
	if scanResults.Sboms != nil {
		sbomRes = &SbomScanResults{
			Results: *scanResults.Sboms,
		}
	}
	var vulRs *VulnerabilityScanResults
	if scanResults.Vulnerabilities != nil {
		vulRs = &VulnerabilityScanResults{
			Results: *scanResults.Vulnerabilities,
		}
	}
	var malwareRes *MalwareScanResults
	if scanResults.Malwares != nil {
		malwareRes = &MalwareScanResults{
			Results: *scanResults.Malwares,
		}
	}
	var secretRes *SecretScanResults
	if scanResults.Secrets != nil {
		secretRes = &SecretScanResults{
			Results: *scanResults.Secrets,
		}
	}
	var rootkitRes *RootkitScanScanResults
	if scanResults.Rootkits != nil {
		rootkitRes = &RootkitScanScanResults{
			Results: *scanResults.Rootkits,
		}
	}
	var misconfigRes *MisconfigurationScanResults
	if scanResults.Misconfigurations != nil {
		misconfigRes = &MisconfigurationScanResults{
			Results: *scanResults.Misconfigurations,
		}
	}
	var exploitRes *ExploitScanResults
	if scanResults.Exploits != nil {
		exploitRes = &ExploitScanResults{
			Results: *scanResults.Exploits,
		}
	}
	return &ScanResults{
		ID:               scanResultID,
		Sbom:             sbomRes,
		Vulnerability:    vulRs,
		Malware:          malwareRes,
		Rootkit:          rootkitRes,
		Secret:           secretRes,
		Misconfiguration: misconfigRes,
		Exploit:          exploitRes,
	}
}

func CreateModelScanResultsFromDB(scanResults *ScanResults) *models.ScanResults {
	var sbomRes models.SbomScan
	if scanResults.Sbom != nil {
		sbomRes = scanResults.Sbom.Results
	}
	var vulRes models.VulnerabilityScan
	if scanResults.Vulnerability != nil {
		vulRes = scanResults.Vulnerability.Results
	}
	var malwareRes models.MalwareScan
	if scanResults.Malware != nil {
		malwareRes = scanResults.Malware.Results
	}
	var secretRes models.SecretScan
	if scanResults.Secret != nil {
		secretRes = scanResults.Secret.Results
	}
	var misconfigRes models.MisconfigurationScan
	if scanResults.Misconfiguration != nil {
		misconfigRes = scanResults.Misconfiguration.Results
	}
	var rootkitRes models.RootkitScan
	if scanResults.Rootkit != nil {
		rootkitRes = scanResults.Rootkit.Results
	}
	var exploitRes models.ExploitScan
	if scanResults.Exploit != nil {
		exploitRes = scanResults.Exploit.Results
	}
	return &models.ScanResults{
		Id:                &scanResults.ID,
		Sboms:             &sbomRes,
		Vulnerabilities:   &vulRes,
		Malwares:          &malwareRes,
		Rootkits:          &rootkitRes,
		Secrets:           &secretRes,
		Misconfigurations: &misconfigRes,
		Exploits:          &exploitRes,
	}
}

func CreateScanResultsSummary(scanResults *ScanResults) *ScanResultsSummary {
	var packagesCount, vulnerabilitiesCount, malwareCount, secretCount, rootkitCount, misconfigurationCount, exploitsCount int
	if scanResults.Sbom != nil {
		packagesCount = len(*scanResults.Sbom.Results.Packages)
	}
	if scanResults.Vulnerability != nil {
		vulnerabilitiesCount = len(*scanResults.Vulnerability.Results.Vulnerabilities)
	}
	if scanResults.Malware != nil {
		malwareCount = len(*scanResults.Malware.Results.Malwares)
	}
	if scanResults.Secret != nil {
		secretCount = len(*scanResults.Secret.Results.Secrets)
	}
	if scanResults.Rootkit != nil {
		rootkitCount = len(*scanResults.Rootkit.Results.Rootkits)
	}
	if scanResults.Misconfiguration != nil {
		misconfigurationCount = len(*scanResults.Misconfiguration.Results.Misconfigurations)
	}
	if scanResults.Exploit != nil {
		exploitsCount = len(*scanResults.Exploit.Results.Exploits)
	}
	return &ScanResultsSummary{
		PackagesCount:          packagesCount,
		VulnerabilitiesCount:   vulnerabilitiesCount,
		MalwaresCount:          malwareCount,
		SecretsCount:           secretCount,
		RootkitsCount:          rootkitCount,
		MisconfigurationsCount: misconfigurationCount,
		ExploitsCount:          exploitsCount,
	}
}

func CreateModelScanResultsSummaryFromDB(summary ScanResultsSummary) *models.ScanResultsSummary {
	return &models.ScanResultsSummary{
		PackagesCount:          &summary.PackagesCount,
		VulnerabilitiesCount:   &summary.VulnerabilitiesCount,
		MalwaresCount:          &summary.MalwaresCount,
		SecretsCount:           &summary.SecretsCount,
		RootkitsCount:          &summary.RootkitsCount,
		MisconfigurationsCount: &summary.MisconfigurationsCount,
		ExploitsCount:          &summary.ExploitsCount,
	}
}

func generateScanResultsID() string {
	return uuid.NewString()
}
