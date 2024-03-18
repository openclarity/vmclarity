package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/scanner/types"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type scanStore struct {
	repo *repo[scanModel]
}

func (h *handler) Scans() types.ScanStore {
	return &scanStore{
		repo: newRepo(h.db, scanModel{}),
	}
}

func (s *scanStore) GetAll(req types.GetScansRequest) ([]types.Scan, error) {
	var filters []interface{}
	if req.State != nil && *req.State != "" {
		filters = append(filters, "state = ?", *req.State)
	}

	var scanModels []scanModel
	if err := s.repo.GetAll(getParams{filters: filters}, &scanModels); err != nil {
		return []types.Scan{}, fmt.Errorf("failed to fetch scans: %w", err)
	}

	// convert
	scans := make([]types.Scan, len(scanModels))
	for idx, scanModel := range scanModels {
		scans[idx] = scanModel.toScan()
	}

	return scans, nil
}

func (s *scanStore) Get(scanID string) (types.Scan, error) {
	var scan scanModel
	if err := s.repo.Get(&scanModel{Id: &scanID}, &scan); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Scan{}, types.ErrNotFound
		}
		return types.Scan{}, fmt.Errorf("failed to get scan: %w", err)
	}

	return scan.toScan(), nil
}

func (s *scanStore) Create(scan types.Scan) (types.Scan, error) {
	if scan.Id != nil && *scan.Id != "" {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "cannot specify id field when creating a new scan",
		}
	}
	if len(scan.Inputs) == 0 {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "cannot have empty inputs field when creating a new scan",
		}
	}

	// Create scan ID
	scanID := uuid.New().String()

	// Create scan
	var toCreate scanModel
	toCreate.fromScan(scan)
	toCreate.Id = &scanID

	if err := s.repo.Create(&toCreate); err != nil {
		return types.Scan{}, fmt.Errorf("failed to create scan: %w", err)
	}

	return s.Get(scanID)
}

func (s *scanStore) Update(scanID string, scan types.Scan) (types.Scan, error) {
	if scan.Id != nil && *scan.Id != scanID {
		return types.Scan{}, &types.PreconditionFailedError{
			Reason: "cannot have different ID when updating a scan",
		}
	}

	// Update scan
	var toUpdate scanModel
	toUpdate.fromScan(scan)
	toUpdate.Id = &scanID

	if err := s.repo.Update(&scanModel{Id: &scanID}, &toUpdate); err != nil {
		return types.Scan{}, fmt.Errorf("failed to update scan: %w", err)
	}

	return s.Get(scanID)
}

func (s *scanStore) Delete(scanID string) error {
	if err := s.repo.Delete(&scanModel{Id: &scanID}); err != nil {
		return fmt.Errorf("failed to delete scan: %w", err)
	}
	return nil
}

type scanModel struct {
	Id             *string           `json:"id,omitempty" gorm:"type:uuid; primaryKey"`
	Status         *types.ScanStatus `json:"status,omitempty" gorm:"embedded"`
	JobsCompleted  *int              `json:"jobsCompleted,omitempty"`
	JobsLeftToRun  *int              `json:"jobsLeftToRun,omitempty"`
	StartTime      *time.Time        `json:"startTime,omitempty"`
	EndTime        *time.Time        `json:"endTime,omitempty"`
	Inputs         datatypes.JSON    `json:"inputs"`
	TimeoutSeconds *int              `json:"timeoutSeconds,omitempty"`
}

func (model *scanModel) toScan() types.Scan {
	var inputs []types.ScanInput
	_ = json.Unmarshal(model.Inputs, &inputs)

	return types.Scan{
		EndTime:        model.EndTime,
		Id:             model.Id,
		Inputs:         inputs,
		JobsCompleted:  model.JobsCompleted,
		JobsLeftToRun:  model.JobsLeftToRun,
		StartTime:      model.StartTime,
		Status:         model.Status,
		TimeoutSeconds: model.TimeoutSeconds,
	}
}

func (model *scanModel) fromScan(scan types.Scan) {
	inputRaw, _ := json.Marshal(scan.Inputs)
	*model = scanModel{
		Id:             scan.Id,
		Status:         scan.Status,
		JobsCompleted:  scan.JobsCompleted,
		JobsLeftToRun:  scan.JobsLeftToRun,
		StartTime:      scan.StartTime,
		EndTime:        scan.EndTime,
		Inputs:         inputRaw,
		TimeoutSeconds: scan.TimeoutSeconds,
	}
}
