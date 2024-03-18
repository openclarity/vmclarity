package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type findingsStore struct {
	repo *repo[findingModel]
}

func (h *handler) ScanFindings() types.ScanFindingStore {
	return &findingsStore{
		repo: newRepo(h.db, findingModel{}),
	}
}

func (s *findingsStore) GetAll(req types.GetScanFindingsRequest) ([]types.ScanFinding, error) {
	var filters []interface{}
	if req.ScanID != nil && *req.ScanID != "" {
		filters = append(filters, "scanID = ?", *req.ScanID)
	}

	var findingModels []findingModel
	if err := s.repo.GetAll(getParams{filters: filters}, &findingModels); err != nil {
		return []types.ScanFinding{}, fmt.Errorf("failed to fetch scan findings: %w", err)
	}

	// convert
	findings := make([]types.ScanFinding, len(findingModels))
	for idx, findingModel := range findingModels {
		findings[idx] = findingModel.toFinding()
	}

	return findings, nil
}

func (s *findingsStore) Get(findingID string) (types.ScanFinding, error) {
	var finding findingModel
	if err := s.repo.Get(&findingModel{Id: &findingID}, &finding); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.ScanFinding{}, types.ErrNotFound
		}
		return types.ScanFinding{}, fmt.Errorf("failed to get scan finding: %w", err)
	}
	return finding.toFinding(), nil
}

func (s *findingsStore) CreateMany(findings ...types.ScanFinding) error {
	toCreate := make([]*findingModel, len(findings))
	for idx, finding := range findings {
		if finding.Id != nil && *finding.Id != "" {
			return &types.PreconditionFailedError{
				Reason: "cannot specify id field when creating a new scan finding",
			}
		}
		if finding.ScanID == nil || *finding.ScanID == "" {
			return &types.PreconditionFailedError{
				Reason: "cannot specify empty scan ID field when creating a new scan finding",
			}
		}

		// Create finding ID
		findingID := uuid.New().String()

		// Create finding
		toCreate[idx] = &findingModel{}
		toCreate[idx].fromFinding(finding)
		toCreate[idx].Id = &findingID
	}

	if err := s.repo.CreateMany(&toCreate); err != nil {
		return fmt.Errorf("failed to create scan findings: %w", err)
	}

	return nil
}

func (s *findingsStore) Delete(req types.DeleteScanFindingsRequest) error {
	if err := s.repo.Delete(&findingModel{
		Id:     req.ID,
		ScanID: req.ScanID,
	}); err != nil {
		return fmt.Errorf("failed to delete scan finding: %w", err)
	}
	return nil
}

type findingModel struct {
	Id          *string               `json:"id,omitempty" gorm:"type:uuid; primaryKey"`
	ScanID      *string               `json:"scanID,omitempty"`
	Annotations datatypes.JSON        `json:"annotations,omitempty"`
	Input       types.ScanInput       `json:"input" gorm:"embedded"`
	FindingInfo datatypes.JSON        `json:"findingInfo"`
	Summary     *types.FindingSummary `json:"summary,omitempty" gorm:"embedded"`
}

func (model *findingModel) toFinding() types.ScanFinding {
	var annotations *types.Annotations
	_ = json.Unmarshal(model.Annotations, &annotations)

	var findingInfo types.ScanFinding_FindingInfo
	_ = findingInfo.UnmarshalJSON(model.FindingInfo)

	return types.ScanFinding{
		Annotations: annotations,
		FindingInfo: findingInfo,
		Id:          model.Id,
		Input:       model.Input,
		ScanID:      model.ScanID,
		Summary:     model.Summary,
	}
}

func (model *findingModel) fromFinding(finding types.ScanFinding) {
	annotationsRaw, _ := json.Marshal(finding.Annotations)
	findingRaw, _ := finding.FindingInfo.MarshalJSON()

	*model = findingModel{
		Id:          finding.Id,
		ScanID:      finding.ScanID,
		Annotations: annotationsRaw,
		Input:       finding.Input,
		FindingInfo: findingRaw,
		Summary:     finding.Summary,
	}
}
