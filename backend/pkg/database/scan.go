package database

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	scansTableName = "scans"
)

type Scan struct {
	Base

	ScanStartTime *time.Time `json:"scan_start_time,omitempty" gorm:"column:scan_start_time"`
	ScanEndTime   *time.Time `json:"scan_end_time,omitempty" gorm:"column:scan_end_time"`

	// ScanConfigId The ID of the config that this scan was initiated from (optionanl)
	ScanConfigId *string `json:"scan_config_id,omitempty" gorm:"column:scan_config_id"`
	// ScanFamiliesConfig The configuration of the scanner families within a scan config
	ScanFamiliesConfig []byte `json:"scan_families_config,omitempty" gorm:"column:scan_families_config"`

	// TargetIDs List of target IDs that are targeted for scanning as part of this scan
	TargetIDs []byte `json:"target_ids,omitempty" gorm:"column:target_ids"`
}

type GetScansParams struct {
	// Filter Odata filter
	Filter *string
	// Page Page number of the query
	Page int
	// PageSize Maximum items to return
	PageSize int
}

type ScansTable interface {
	GetScansAndTotal(params GetScansParams) ([]*Scan, int64, error)
	GetScan(scanID string) (*Scan, error)
	CheckExist(scanConfigID string) (*Scan, bool, error)
	UpdateScan(scan *Scan, scanID string) (*Scan, error)
	SaveScan(scan *Scan, scanID string) (*Scan, error)
	DeleteScan(scanID string) error
	CreateScan(scan *Scan) (*Scan, error)
}

type ScansTableHandler struct {
	scansTable *gorm.DB
}

func (db *Handler) ScansTable() ScansTable {
	return &ScansTableHandler{
		scansTable: db.DB.Table(scansTableName),
	}
}

func (s *ScansTableHandler) CheckExist(scanConfigID string) (*Scan, bool, error) {
	var scan *Scan

	if err := s.scansTable.Where("scan_config_id = ?", scanConfigID).First(&scan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return scan, true, nil
}

func (s *ScansTableHandler) GetScansAndTotal(params GetScansParams) ([]*Scan, int64, error) {
	var count int64
	var scans []*Scan

	tx := s.scansTable

	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count total: %w", err)
	}

	if err := tx.Find(&scans).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find scans: %w", err)
	}

	return scans, count, nil
}

func (s *ScansTableHandler) CreateScan(scan *Scan) (*Scan, error) {
	if err := s.scansTable.Create(scan).Error; err != nil {
		return nil, err
	}
	return scan, nil
}

func (s *ScansTableHandler) SaveScan(scan *Scan, scanID string) (*Scan, error) {
	var err error
	scan.ID, err = uuid.FromString(scanID)
	if err != nil {
		return nil, err
	}

	if err := s.scansTable.Save(scan).Error; err != nil {
		return nil, err
	}

	return scan, nil
}

func (s *ScansTableHandler) UpdateScan(scan *Scan, scanID string) (*Scan, error) {
	var err error
	scan.ID, err = uuid.FromString(scanID)
	if err != nil {
		return nil, err
	}

	selectClause := []string{}
	if len(scan.ScanFamiliesConfig) > 0 {
		selectClause = append(selectClause, "scan_families_config")
	}
	if scan.ScanConfigId != nil {
		selectClause = append(selectClause, "scan_config_id")
	}
	if scan.ScanStartTime != nil {
		selectClause = append(selectClause, "scan_start_time")
	}
	if scan.ScanEndTime != nil {
		selectClause = append(selectClause, "scan_end_time")
	}
	if scan.TargetIDs != nil {
		selectClause = append(selectClause, "target_ids")
	}


	if err := s.scansTable.Model(scan).Select(selectClause).Updates(scan).Error ; err != nil {
		return nil, err
	}
	return scan, nil
}

func (s *ScansTableHandler) GetScan(scanID string) (*Scan, error) {
	var scan *Scan

	if err := s.scansTable.Where("id = ?", scanID).First(&scan).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan by id %q: %w", scanID, err)
	}

	return scan, nil
}

func (s *ScansTableHandler) DeleteScan(scanID string) error {
	if err := s.scansTable.Delete(&Scan{}, scanID).Error; err != nil {
		return fmt.Errorf("failed to delete scan: %w", err)
	}
	return nil
}
