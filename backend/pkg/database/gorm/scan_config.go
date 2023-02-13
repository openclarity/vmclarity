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

package gorm

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/common"
	"github.com/openclarity/vmclarity/backend/pkg/database/types"
)

const (
	scanConfigsTableName = "scan_configs"
)

type ScanConfig struct {
	Base

	Name *string `json:"name,omitempty" gorm:"column:name"`

	// ScanFamiliesConfig The configuration of the scanner families within a scan config
	ScanFamiliesConfig []byte `json:"scan_families_config,omitempty" gorm:"column:scan_families_config"`
	Scheduled          []byte `json:"scheduled,omitempty" gorm:"column:scheduled"`
	Scope              []byte `json:"scope,omitempty" gorm:"column:scope"`
}

type ScanConfigsTableHandler struct {
	scanConfigsTable *gorm.DB
}

func (db *Handler) ScanConfigsTable() types.ScanConfigsTable {
	return &ScanConfigsTableHandler{
		scanConfigsTable: db.DB.Table(scanConfigsTableName),
	}
}

func (s *ScanConfigsTableHandler) GetScanConfigs(params models.GetScanConfigsParams) (models.ScanConfigs, error) {
	var scanConfigs []ScanConfig

	tx := s.scanConfigsTable

	if err := tx.Find(&scanConfigs).Error; err != nil {
		return models.ScanConfigs{}, fmt.Errorf("failed to find scan configs: %w", err)
	}

	converted, err := ConvertToRestScanConfigs(scanConfigs)
	if err != nil {
		return models.ScanConfigs{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
	}
	return converted, nil
}

func (s *ScanConfigsTableHandler) GetScanConfig(scanConfigID models.ScanConfigID) (models.ScanConfig, error) {
	var dbScanConfig ScanConfig
	if err := s.scanConfigsTable.Where("id = ?", scanConfigID).First(&dbScanConfig).Error; err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to get scan config by id %q: %w", scanConfigID, err)
	}

	converted, err := ConvertToRestScanConfig(dbScanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
	}
	return converted, nil
}

func (s *ScanConfigsTableHandler) CreateScanConfig(scanConfig models.ScanConfig) (models.ScanConfig, error) {
	// check if there is already a scan config with that name.
	existingSR, exist, err := s.checkExist(*scanConfig.Name)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to check existing scan config: %w", err)
	}
	if exist {
		converted, err := ConvertToRestScanConfig(existingSR)
		if err != nil {
			return models.ScanConfig{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
		}
		return converted, &common.ConflictError{
			Reason: fmt.Sprintf("Scan config exists with name=%s", *existingSR.Name),
		}
	}

	dbScanConfig, err := ConvertToDBScanConfig(scanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert API model to DB model: %w", err)
	}

	if err := s.scanConfigsTable.Create(&dbScanConfig).Error; err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to create scan config in db: %w", err)
	}

	converted, err := ConvertToRestScanConfig(dbScanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
	}
	return converted, nil
}

func (s *ScanConfigsTableHandler) SaveScanConfig(scanConfig models.ScanConfig) (models.ScanConfig, error) {
	dbScanConfig, err := ConvertToDBScanConfig(scanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert API model to DB model: %w", err)
	}

	if err := s.scanConfigsTable.Save(&dbScanConfig).Error; err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to save scan config in db: %w", err)
	}

	converted, err := ConvertToRestScanConfig(dbScanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
	}
	return converted, nil
}

func (s *ScanConfigsTableHandler) UpdateScanConfig(scanConfig models.ScanConfig) (models.ScanConfig, error) {
	dbScanConfig, err := ConvertToDBScanConfig(scanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert API model to DB model: %w", err)
	}

	if err := s.scanConfigsTable.Model(dbScanConfig).Updates(&dbScanConfig).Error; err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to update scan config in db: %w", err)
	}

	converted, err := ConvertToRestScanConfig(dbScanConfig)
	if err != nil {
		return models.ScanConfig{}, fmt.Errorf("failed to convert DB model to API model: %w", err)
	}
	return converted, nil
}

func (s *ScanConfigsTableHandler) DeleteScanConfig(scanConfigID models.ScanConfigID) error {
	if err := s.scanConfigsTable.Where("id = ?", scanConfigID).Delete(&Scan{}).Error; err != nil {
		return fmt.Errorf("failed to delete scan config: %w", err)
	}
	return nil
}

func (s *ScanConfigsTableHandler) checkExist(name string) (ScanConfig, bool, error) {
	var scanConfig ScanConfig

	tx := s.scanConfigsTable.WithContext(context.Background())

	if err := tx.Where("name = ?", name).First(&scanConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ScanConfig{}, false, nil
		}
		return ScanConfig{}, false, fmt.Errorf("failed to query: %w", err)
	}

	return scanConfig, true, nil
}
