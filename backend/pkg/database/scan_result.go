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
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	scanResultsTableName = "scan_results"
)

type ScanResult struct {
	Base

	ScanID   string `json:"scan_id,omitempty" gorm:"column:scan_id"`
	TargetID string `json:"target_id,omitempty" gorm:"column:target_id"`

	Exploits          []byte `json:"exploits,omitempty" gorm:"column:exploits"`
	Malware           []byte `json:"malware,omitempty" gorm:"column:malware"`
	Misconfigurations []byte `json:"misconfigurations,omitempty" gorm:"column:misconfigurations"`
	Rootkits          []byte `json:"rootkits,omitempty" gorm:"column:rootkits"`
	Sboms             []byte `json:"sboms,omitempty" gorm:"column:sboms"`
	Secrets           []byte `json:"secrets,omitempty" gorm:"column:secrets"`
	Status            []byte `json:"status,omitempty" gorm:"column:status"`
	Vulnerabilities   []byte `json:"vulnerabilities,omitempty" gorm:"column:vulnerabilities"`
}

type GetScanResultsParams struct {
	// Filter Odata filter
	Filter *string
	// Select Odata select
	Select *string
	// Page Page number of the query
	Page int
	// PageSize Maximum items to return
	PageSize int
}

type GetScanResultsScanResultIDParams struct {
	// Select Odata select
	Select *string
}

//nolint:interfacebloat
type ScanResultsTable interface {
	CreateScanResult(scanResults *ScanResult) (*ScanResult, error)
	GetScanResultsAndTotal(params GetScanResultsParams) ([]*ScanResult, int64, error)
	CheckExist(scanID string, targetID string) (*ScanResult, bool, error)
	GetScanResult(scanResultID string, params GetScanResultsScanResultIDParams) (*ScanResult, error)
	UpdateScanResult(scanResults *ScanResult, scanResultID string) (*ScanResult, error)
	SaveScanResult(scanResults *ScanResult, scanResultID string) (*ScanResult, error)
}

type ScanResultsTableHandler struct {
	scanResultsTable *gorm.DB
}

func (db *Handler) ScanResultsTable() ScanResultsTable {
	return &ScanResultsTableHandler{
		scanResultsTable: db.DB.Table(scanResultsTableName),
	}
}

func (s *ScanResultsTableHandler) CheckExist(scanID string, targetID string) (*ScanResult, bool, error) {
	var scanResult *ScanResult

	if err := s.scanResultsTable.Where("scan_id = ? AND target_id = ?", scanID, targetID).First(&scanResult).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return scanResult, true, nil
}

func (s *ScanResultsTableHandler) GetScanResult(scanResultID string, params GetScanResultsScanResultIDParams) (*ScanResult, error) {
	var scanResult *ScanResult

	if err := s.scanResultsTable.Where("id = ?", scanResultID).First(&scanResult).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan result by id %q: %w", scanResultID, err)
	}

	return scanResult, nil
}

func (s *ScanResultsTableHandler) CreateScanResult(scanResult *ScanResult) (*ScanResult, error) {
	if err := s.scanResultsTable.Create(scanResult).Error; err != nil {
		return nil, err
	}
	return scanResult, nil
}

func (s *ScanResultsTableHandler) GetScanResultsAndTotal(params GetScanResultsParams) ([]*ScanResult, int64, error) {
	var count int64
	var scanResults []*ScanResult

	tx := s.scanResultsTable

	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count total: %w", err)
	}

	if err := tx.Find(&scanResults).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find scan results: %w", err)
	}

	return scanResults, count, nil
}

func (s *ScanResultsTableHandler) SaveScanResult(scanResult *ScanResult, scanResultID string) (*ScanResult, error) {
	var err error
	scanResult.ID, err = uuid.FromString(scanResultID)
	if err != nil {
		return nil, err
	}
	if err := s.scanResultsTable.Save(scanResult).Error; err != nil {
		return nil, err
	}

	return scanResult, nil
}

func (s *ScanResultsTableHandler) UpdateScanResult(scanResult *ScanResult, scanResultID string) (*ScanResult, error) {
	var err error
	scanResult.ID, err = uuid.FromString(scanResultID)
	if err != nil {
		return nil, err
	}

	selectClause := []string{}
	if len(scanResult.ScanID) > 0 {
		selectClause = append(selectClause, "scan_id")
	}
	if len(scanResult.TargetID) > 0 {
		selectClause = append(selectClause, "target_id")
	}
	if scanResult.Sboms != nil {
		selectClause = append(selectClause, "sboms")
	}
	if scanResult.Status != nil {
		selectClause = append(selectClause, "status")
	}
	if scanResult.Rootkits != nil {
		selectClause = append(selectClause, "rootkits")
	}
	if scanResult.Malware != nil {
		selectClause = append(selectClause, "malware")
	}
	if scanResult.Vulnerabilities != nil {
		selectClause = append(selectClause, "vulnerabilities")
	}
	if scanResult.Secrets != nil {
		selectClause = append(selectClause, "secrets")
	}
	if scanResult.Misconfigurations != nil {
		selectClause = append(selectClause, "misconfigurations")
	}
	if scanResult.Exploits != nil {
		selectClause = append(selectClause, "exploits")
	}

	if err := s.scanResultsTable.Model(scanResult).Select(selectClause).Updates(scanResult).Error ; err != nil {
		return nil, err
	}

	return scanResult, nil
}
