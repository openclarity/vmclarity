package local

import (
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/gorm"
)

type scanResultModel struct {
	baseModel
	Data types.ScanResult
}

type scanResultStore struct {
	db *gorm.DB
}

func (h *handler) ScanResultStore() types.ScanResultStore {
	return &scanResultStore{
		db: h.db,
	}
}

func (s scanResultStore) GetScanResult(scanID types.ScanID) (types.ScanResult, error) {
	var scanResult scanResultModel
	if err := s.db.First(&scanResult, scanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ScanResult{}, types.ErrNotFound
		}
		return types.ScanResult{}, fmt.Errorf("failed to get scan result: %w", err)
	}
	return scanResult.Data, nil
}

func (s scanResultStore) CreateScanResult(result types.ScanResult) (types.ScanResult, error) {
	if result.ScanID != "" {
		return types.ScanResult{}, &types.PreconditionFailedError{
			Reason: "can not specify empty scanID field when creating a new scan result",
		}
	}

	// Create scan
	var toCreate scanResultModel
	toCreate.ID = result.ScanID
	toCreate.Data = result
	if err := s.db.Create(&toCreate).Error; err != nil {
		return types.ScanResult{}, fmt.Errorf("failed to create scan result: %w", err)
	}

	return s.GetScanResult(result.ScanID)
}

func (s scanResultStore) UpdateScanResult(scanID types.ScanID, result types.ScanResult) (types.ScanResult, error) {
	if result.ScanID != scanID {
		return types.ScanResult{}, &types.PreconditionFailedError{
			Reason: "can not have different ID when updating a scan result",
		}
	}

	// Create scan
	var toUpdate scanResultModel
	toUpdate.ID = scanID
	toUpdate.Data = result
	if err := s.db.Save(&toUpdate).Error; err != nil {
		return types.ScanResult{}, fmt.Errorf("failed to update scan: %w", err)
	}

	return s.GetScanResult(scanID)
}
