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

type ScanResultsInt interface{}

// TODO after db design.
type ScanResults struct {
	ID               string
	Sbom             SbomScanResults
	Vulnerability    VulnerabilityScanResults
	Malware          MalwareScanResults
	Rootkit          RootkitScanScanResults
	Secret           SecretScanResults
	Misconfiguration MisconfigurationScanResults
	Exploit          ExploitScanResults
}

// TODO after db design.
type SbomScanResults struct {
	ID string
}

type VulnerabilityScanResults struct {
	ID string
}

type MalwareScanResults struct {
	ID string
}

type RootkitScanScanResults struct {
	ID string
}

type SecretScanResults struct {
	ID string
}

type MisconfigurationScanResults struct {
	ID string
}

type ExploitScanResults struct {
	ID string
}

//go:generate $GOPATH/bin/mockgen --build_flags=--mod=mod -destination=./mock_scan_results.go -package=database github.com/openclarity/vmclarity/backend/pkg/database ScanResultsTable
type ScanResultsTable interface {
	List(targetID models.TargetID, params models.GetTargetsTargetIDScanresultsParams) ([]models.ScanResults, error)
	Create(targetID models.TargetID, scanResults *ScanResults) (models.ScanResultsSummary, error)
	Get(targetID models.TargetID, scanID models.ScanID, params models.GetTargetsTargetIDScanresultsScanIDParams) (ScanResultsInt, error)
	Update(targetID models.TargetID, scanID models.ScanID, scanResults *ScanResults) (models.ScanResultsSummary, error)
}

type ScanResultsTableHandler struct {
	db *gorm.DB
}

func (db *Handler) ScanResultsTable() ScanResultsTable {
	return &ScanResultsTableHandler{
		db: db.DB,
	}
}

func (s *ScanResultsTableHandler) List(targetID models.TargetID, params models.GetTargetsTargetIDScanresultsParams,
) ([]models.ScanResults, error) {
	return []models.ScanResults{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) Create(targetID models.TargetID, scanResults *ScanResults,
) (models.ScanResultsSummary, error) {
	return models.ScanResultsSummary{}, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) Get(
	targetID models.TargetID,
	scanID models.ScanID,
	params models.GetTargetsTargetIDScanresultsScanIDParams,
) (ScanResultsInt, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ScanResultsTableHandler) Update(
	targetID models.TargetID,
	scanID models.ScanID,
	scanResults *ScanResults,
) (models.ScanResultsSummary, error) {
	return models.ScanResultsSummary{}, fmt.Errorf("not implemented")
}

// TODO after db design.
func CreateScanResults(scanResults *models.ScanResults) *ScanResults {
	return &ScanResults{
		ID: *scanResults.Id,
	}
}
