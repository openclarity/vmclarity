package database

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

const (
	scanConfigsTableName = "scan_configs"
)

type ScanConfig struct {
	gorm.Model

	Name *string `json:"name,omitempty" gorm:"column:name"`

	// ScanFamiliesConfig The configuration of the scanner families within a scan config
	ScanFamiliesConfig []byte `json:"scan_families_config,omitempty" gorm:"column:scan_families_config"`
	Scheduled          []byte `json:"scheduled,omitempty" gorm:"column:scheduled"`
	Scope              []byte `json:"scope,omitempty" gorm:"column:scope"`
}

type GetScanConfigsParams struct {
	// Filter Odata filter
	Filter *string
	// Page Page number of the query
	Page int
	// PageSize Maximum items to return
	PageSize int
}

type ScanConfigsTable interface {
	GetScanConfigsAndTotal(params GetScanConfigsParams) ([]*ScanConfig, int64, error)
	GetScanConfig(scanConfigID string) (*ScanConfig, error)
	CheckExist(name string) (*ScanConfig, bool, error)
	UpdateScanConfig(scanConfig *ScanConfig, scanConfigID string) (*ScanConfig, error)
	SaveScanConfig(scanConfig *ScanConfig, scanConfigID string) (*ScanConfig, error)
	DeleteScanConfig(scanConfigID string) error
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

func (s *ScanConfigsTableHandler) CheckExist(name string) (*ScanConfig, bool, error) {
	var scanConfig *ScanConfig

	if err := s.scanConfigsTable.Where("name = ?", name).First(&scanConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return scanConfig, true, nil
}

func (s *ScanConfigsTableHandler) GetScanConfigsAndTotal(params GetScanConfigsParams) ([]*ScanConfig, int64, error) {
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

func (s *ScanConfigsTableHandler) SaveScanConfig(scanConfig *ScanConfig, scanConfigID string) (*ScanConfig, error) {
	id, err := strconv.Atoi(scanConfigID)
	if err != nil {
		return nil, err
	}
	scanConfig.ID = uint(id)

	if err := s.scanConfigsTable.Save(scanConfig).Error ; err != nil {
		return nil, err
	}

	return scanConfig, err
}

func (s *ScanConfigsTableHandler) UpdateScanConfig(scanConfig *ScanConfig, scanConfigID string) (*ScanConfig, error) {
	id, err := strconv.Atoi(scanConfigID)
	if err != nil {
		return nil, err
	}
	scanConfig.ID = uint(id)

	selectClause := []string{}
	if len(scanConfig.ScanFamiliesConfig) > 0 {
		selectClause = append(selectClause, "scan_families_config")
	}
	if scanConfig.Name != nil {
		selectClause = append(selectClause, "name")
	}
	if len(scanConfig.Scheduled) > 0 {
		selectClause = append(selectClause, "scheduled")
	}
	if len(scanConfig.Scope) > 0 {
		selectClause = append(selectClause, "scope")
	}

	if err := s.scanConfigsTable.Model(scanConfig).Select(selectClause).Updates(scanConfig).Error ; err != nil {
		return nil, err
	}

	return scanConfig, err
}

func (s *ScanConfigsTableHandler) GetScanConfig(scanConfigID string) (*ScanConfig, error) {
	var scanConfig *ScanConfig

	if err := s.scanConfigsTable.Where("id = ?", scanConfigID).First(&scanConfig).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan config by id %q: %w", scanConfigID, err)
	}

	return scanConfig, nil
}

func (s *ScanConfigsTableHandler) DeleteScanConfig(scanConfigID string) error {
	if err := s.scanConfigsTable.Delete(&Scan{}, scanConfigID).Error; err != nil {
		return fmt.Errorf("failed to delete scan config: %w", err)
	}
	return nil
}
