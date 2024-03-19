package local

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"github.com/openclarity/vmclarity/scanner/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type findingsStore struct {
	repo *repo[findingModel]
}

func (h *handler) ScanFindings() store.ScanFindingStore {
	return &findingsStore{
		repo: newRepo(h.db, findingModel{}),
	}
}

func (s *findingsStore) GetAll(req store.GetScanFindingsRequest) ([]types.ScanFinding, error) {
	var filters [][]interface{}
	if req.ScanID != nil && *req.ScanID != "" {
		filters = append(filters, []interface{}{"scan_id = ?", *req.ScanID})
	}

	// TODO(ramizpolic): This meta filter is not applied properly, why?
	if req.MetaSelector != nil && len(*req.MetaSelector) > 0 {
		selector := getMetaKVSelectors("annotations", *req.MetaSelector)
		if selector != nil {
			filters = append(filters, []interface{}{*selector})
		}
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
			return types.ScanFinding{}, store.ErrNotFound
		}
		return types.ScanFinding{}, fmt.Errorf("failed to get scan finding: %w", err)
	}
	return finding.toFinding(), nil
}

func (s *findingsStore) CreateMany(scanID string, findings ...types.ScanFinding) ([]types.ScanFinding, error) {
	toCreate := make([]*findingModel, len(findings))
	for idx, finding := range findings {
		if finding.Id != nil && *finding.Id != "" {
			return nil, &store.PreconditionFailedError{
				Reason: "cannot specify id field when creating a new scan finding",
			}
		}
		if finding.ScanID != nil && *finding.ScanID != scanID {
			return nil, &store.PreconditionFailedError{
				Reason: "cannot specify different scan ID field when creating a new scan finding",
			}
		}

		// Create finding ID
		findingID := uuid.New().String()

		// Create finding
		toCreate[idx] = &findingModel{}
		toCreate[idx].fromFinding(finding)
		toCreate[idx].Id = &findingID
		toCreate[idx].ScanID = &scanID
	}

	if err := s.repo.CreateMany(&toCreate); err != nil {
		return nil, fmt.Errorf("failed to create scan findings: %w", err)
	}

	return s.GetAll(store.GetScanFindingsRequest{
		ScanID: &scanID,
	})
}

func (s *findingsStore) Delete(req store.DeleteScanFindingsRequest) error {
	if err := s.repo.Delete(&findingModel{
		ScanID: req.ScanID,
	}); err != nil {
		return fmt.Errorf("failed to delete scan finding: %w", err)
	}
	return nil
}

type findingModel struct {
	Id          *string                                            `json:"id,omitempty" gorm:"type:uuid; primaryKey"`
	ScanID      *string                                            `json:"scanId,omitempty" gorm:"type:uuid;"`
	Annotations *datatypes.JSONType[map[string]string]             `json:"annotations,omitempty"`
	Input       *types.ScanInput                                   `json:"input,omitempty" gorm:"embedded"`
	FindingInfo *datatypes.JSONType[types.ScanFinding_FindingInfo] `json:"findingInfo,omitempty"`
	Summary     *types.FindingSummary                              `json:"summary,omitempty" gorm:"embedded"`
}

func (findingModel) TableName() string { return "findings" }

func (model *findingModel) toFinding() types.ScanFinding {
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

	var findingInfo *types.ScanFinding_FindingInfo
	if model.FindingInfo != nil {
		findingInfo = new(types.ScanFinding_FindingInfo)
		*findingInfo = model.FindingInfo.Data()
	}

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
	var annotations *datatypes.JSONType[map[string]string]
	if finding.Annotations != nil {
		annotations = new(datatypes.JSONType[map[string]string])
		*annotations = datatypes.NewJSONType(types.AnnotationsAsMap(finding.Annotations))
	}

	var findingInfo *datatypes.JSONType[types.ScanFinding_FindingInfo]
	if finding.FindingInfo != nil {
		findingInfo = new(datatypes.JSONType[types.ScanFinding_FindingInfo])
		*findingInfo = datatypes.NewJSONType(*finding.FindingInfo)
	}

	*model = findingModel{
		Id:          finding.Id,
		ScanID:      finding.ScanID,
		Annotations: annotations,
		Input:       finding.Input,
		FindingInfo: findingInfo,
		Summary:     finding.Summary,
	}
}
