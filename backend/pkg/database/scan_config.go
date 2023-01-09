package database

import (
	"fmt"
	"strconv"

	"github.com/openclarity/vmclarity/api/models"
	"gorm.io/gorm"
)

const (
	scanConfigsTableName = "scan_configs"
)

type ScanConfig struct {
	gorm.Model

	Name string `json:"name,omitempty" gorm:"column:name"`

	// ScanFamiliesConfig The configuration of the scanner families within a scan config
	ScanFamiliesConfig []byte `json:"scan_families_config,omitempty" gorm:"column:scan_families_config"`
	Scheduled          []byte `json:"scheduled,omitempty" gorm:"column:scheduled"`
	Scope              []byte `json:"scope,omitempty" gorm:"column:scope"`
}

type ScanConfigsTable interface {
	GetScanConfigsAndTotal(params models.GetScanConfigsParams) ([]*ScanConfig, int64, error)
	GetScanConfig(scanConfigID models.ScanConfigID) (*ScanConfig, error)
	UpdateScanConfig(scanConfig *ScanConfig, scanConfigID models.ScanConfigID) (*ScanConfig, error)
	DeleteScanConfig(scanConfigID models.ScanConfigID) error
	CreateScanConfig(scanConfig *ScanConfig) (*ScanConfig, error)
}

type ScanConfigsTableHandler struct {
	scanConfigsTable *gorm.DB
}

func (db *Handler) ScanConfigsTable() ScanConfigsTable {
	return &ScanConfigsTableHandler{
		scanConfigsTable: db.DB.Table(scanConfigsTableName),
	}
}

func (s *ScanConfigsTableHandler) GetScanConfigsAndTotal(params models.GetScanConfigsParams) ([]*ScanConfig, int64, error) {
	var count int64
	var scanConfigs []*ScanConfig

	tx := s.scanConfigsTable

	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count total: %w", err)
	}

	if err := tx.Find(&scanConfigs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find scan configs: %w", err)
	}

	return scanConfigs, count, nil
}

func (s *ScanConfigsTableHandler) CreateScanConfig(scanConfig *ScanConfig) (*ScanConfig, error) {
	if err := s.scanConfigsTable.Create(scanConfig).Error; err != nil {
		return nil, err
	}
	return scanConfig, nil
}

func (s *ScanConfigsTableHandler) UpdateScanConfig(scanConfig *ScanConfig, scanConfigID models.ScanConfigID) (*ScanConfig, error) {
	id, err := strconv.Atoi(scanConfigID)
	if err != nil {
		return nil, err
	}
	scanConfig.ID = uint(id)
	s.scanConfigsTable.Save(scanConfig)

	return scanConfig, err
}

func (s *ScanConfigsTableHandler) GetScanConfig(scanConfigID models.ScanConfigID) (*ScanConfig, error) {
	var scanConfig *ScanConfig

	if err := s.scanConfigsTable.Where("id = ?", scanConfigID).First(&scanConfig).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan config by id %q: %w", scanConfigID, err)
	}

	return scanConfig, nil
}

func (s *ScanConfigsTableHandler) DeleteScanConfig(scanConfigID models.ScanConfigID) error {
	if err := s.scanConfigsTable.Delete(&Scan{}, scanConfigID).Error; err != nil {
		return fmt.Errorf("failed to delete scan config: %w", err)
	}
	return nil
}
