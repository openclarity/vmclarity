package local

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/gorm"
)

type scanModel struct {
	baseModel
	Data types.Scan
}

type scanStore struct {
	db *gorm.DB
}

func (h *handler) ScanStore() types.ScanStore {
	return &scanStore{
		db: h.db,
	}
}

func (s *scanStore) GetScans(params types.GetScansParams) (types.Scans, error) {
	// Get queries
	page := 0
	if params.Page != nil {
		page = *params.Page - 1
	}

	pageSize := 50
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	query := s.db
	if params.State != nil {
		query = s.db.Where("status.state = ?", *params.State)
	}

	// Fetch scans
	var scans []scanModel
	if err := query.Limit(pageSize).Offset(page * pageSize).Find(&scans).Error; err != nil {
		return types.Scans{}, fmt.Errorf("failed to get scans: %w", err)
	}

	// Convert to proper type
	var result types.Scans
	for _, scan := range scans {
		result.Count += 1
		result.Items = append(result.Items, scan.Data)
	}

	return result, nil
}

func (s *scanStore) GetScan(scanID types.ScanID) (types.Scan, error) {
	var scan scanModel
	if err := s.db.First(&scan, scanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Scan{}, types.ErrNotFound
		}
		return types.Scan{}, fmt.Errorf("failed to get scan: %w", err)
	}
	return scan.Data, nil
}

func (s *scanStore) CreateScan(scan types.Scan) (types.Scan, error) {
	// Check the user didn't provide an ID
	if scan.Id != nil && *scan.Id != "" {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "can not specify id field when creating a new scan",
		}
	}
	if len(scan.Inputs) == 0 {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "can not have empty inputs field when creating a new scan",
		}
	}

	// Set scan ID
	scanID := uuid.New().String()
	scan.Id = &scanID

	// Create scan
	var toCreate scanModel
	toCreate.ID = scanID
	toCreate.Data = scan
	if err := s.db.Create(&toCreate).Error; err != nil {
		return types.Scan{}, fmt.Errorf("failed to create scan: %w", err)
	}

	return s.GetScan(scanID)
}

func (s *scanStore) UpdateScan(scanID types.ScanID, scan types.Scan) (types.Scan, error) {
	if scan.Id == nil || *scan.Id != scanID {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "can not have different ID when updating a scan",
		}
	}

	// Create scan
	var toUpdate scanModel
	toUpdate.ID = scanID
	toUpdate.Data = scan
	if err := s.db.Save(&toUpdate).Error; err != nil {
		return types.Scan{}, fmt.Errorf("failed to update scan: %w", err)
	}

	return s.GetScan(scanID)
}

func (s *scanStore) DeleteScan(scanID types.ScanID) error {
	if err := s.db.Delete(&scanModel{}, scanID).Error; err != nil {
		return fmt.Errorf("failed to delete scan: %w", err)
	}
	return nil
}
