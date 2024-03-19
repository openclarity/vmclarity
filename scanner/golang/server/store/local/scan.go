package local

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"github.com/openclarity/vmclarity/scanner/types"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type scanStore struct {
	repo *repo[scanModel]
}

func (h *handler) Scans() store.ScanStore {
	return &scanStore{
		repo: newRepo(h.db.Debug(), scanModel{}),
	}
}

func (s *scanStore) GetAll(req store.GetScansRequest) ([]types.Scan, error) {
	var filters [][]interface{}
	if req.State != nil && *req.State != "" {
		filters = append(filters, []interface{}{"state = ?", *req.State})
	}
	if req.MetaSelector != nil && len(*req.MetaSelector) > 0 {
		selector := getMetaKVSelectors("annotations", *req.MetaSelector)
		if selector != nil {
			filters = append(filters, []interface{}{selector})
		}
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
			return types.Scan{}, store.ErrNotFound
		}
		return types.Scan{}, fmt.Errorf("failed to get scan: %w", err)
	}

	return scan.toScan(), nil
}

func (s *scanStore) Create(scan types.Scan) (types.Scan, error) {
	if scan.Id != nil && *scan.Id != "" {
		return types.Scan{}, &store.PreconditionFailedError{
			Reason: "cannot specify id field when creating a new scan",
		}
	}
	if scan.Inputs == nil || len(*scan.Inputs) == 0 {
		return types.Scan{}, &store.PreconditionFailedError{
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
		return types.Scan{}, &store.PreconditionFailedError{
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
	Annotations              *datatypes.JSONType[map[string]string] `json:"annotations,omitempty"`
	Id                       *string                                `json:"id,omitempty" gorm:"type:uuid; primaryKey"`
	Status                   *types.ScanStatus                      `json:"status,omitempty" gorm:"embedded"`
	Summary                  *types.ScanSummary                     `json:"summary,omitempty" gorm:"embedded"`
	Scanner                  *datatypes.JSONType[types.ScannerInfo] `json:"scanner,omitempty"`
	SubmitTime               *time.Time                             `json:"submitTime,omitempty"`
	StartTime                *time.Time                             `json:"startTime,omitempty"`
	EndTime                  *time.Time                             `json:"endTime,omitempty"`
	Inputs                   *datatypes.JSONType[[]types.ScanInput] `json:"inputs,omitempty"`
	InProgressTimeoutSeconds *int                                   `json:"inProgressTimeoutSeconds,omitempty"`
	PendingTimeoutSeconds    *int                                   `json:"pendingTimeoutSeconds,omitempty"`
}

func (scanModel) TableName() string { return "scans" }

func (model *scanModel) toScan() types.Scan {
	var annotations *types.Annotations
	if model.Annotations != nil {
		annotations = new(types.Annotations)
		for key, value := range model.Annotations.Data() {
			key, value := key, value // loop error
			*annotations = append(*annotations, types.Annotation{
				Key:   &key,
				Value: &value,
			})
		}
	}

	var scannerInfo *types.ScannerInfo
	if model.Scanner != nil {
		scannerInfo = new(types.ScannerInfo)
		*scannerInfo = model.Scanner.Data()
	}

	var inputs *[]types.ScanInput
	if model.Inputs != nil {
		inputs = new([]types.ScanInput)
		*inputs = model.Inputs.Data()
	}

	return types.Scan{
		Annotations:              annotations,
		EndTime:                  model.EndTime,
		Id:                       model.Id,
		InProgressTimeoutSeconds: model.InProgressTimeoutSeconds,
		Inputs:                   inputs,
		PendingTimeoutSeconds:    model.PendingTimeoutSeconds,
		Scanner:                  scannerInfo,
		StartTime:                model.StartTime,
		Status:                   model.Status,
		SubmitTime:               model.SubmitTime,
		Summary:                  model.Summary,
	}
}

func (model *scanModel) fromScan(scan types.Scan) {
	var inputs *datatypes.JSONType[[]types.ScanInput]
	if scan.Inputs != nil {
		inputs = new(datatypes.JSONType[[]types.ScanInput])
		*inputs = datatypes.NewJSONType(*scan.Inputs)
	}

	var annotations *datatypes.JSONType[map[string]string]
	if scan.Annotations != nil {
		annotations = new(datatypes.JSONType[map[string]string])
		*annotations = datatypes.NewJSONType(types.AnnotationsAsMap(scan.Annotations))
	}

	var scanner *datatypes.JSONType[types.ScannerInfo]
	if scan.Scanner != nil {
		scanner = new(datatypes.JSONType[types.ScannerInfo])
		*scanner = datatypes.NewJSONType(*scan.Scanner)
	}

	*model = scanModel{
		Annotations:              annotations,
		Id:                       scan.Id,
		Status:                   scan.Status,
		Summary:                  scan.Summary,
		Scanner:                  scanner,
		SubmitTime:               scan.SubmitTime,
		StartTime:                scan.StartTime,
		EndTime:                  scan.EndTime,
		Inputs:                   inputs,
		InProgressTimeoutSeconds: scan.InProgressTimeoutSeconds,
		PendingTimeoutSeconds:    scan.PendingTimeoutSeconds,
	}
}
