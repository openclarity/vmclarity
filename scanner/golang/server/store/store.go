package store

import (
	"errors"
	"github.com/openclarity/vmclarity/scanner/types"
	"strings"
)

const MetaSelectorSeparator = "="

var ErrNotFound = errors.New("not found")

type PreconditionFailedError struct {
	Reason string
}

func (e *PreconditionFailedError) Error() string {
	return "Precondition failed: " + e.Reason
}

type GetScansRequest struct {
	State        *string   `json:"state,omitempty"`
	MetaSelector *[]string `json:"metaSelector,omitempty"`
}

func (r *GetScansRequest) GetMetaKVSelectors() map[string]string {
	if r.MetaSelector == nil {
		return nil
	}

	// this extracts selectors such as key=value specified for metaSelector query param
	kvSelectors := make(map[string]string)
	for _, item := range *r.MetaSelector {
		items := strings.SplitN(item, MetaSelectorSeparator, 2)
		if len(items) < 2 {
			continue
		}
		kvSelectors[items[0]] = items[1]
	}
	return kvSelectors
}

type ScanStore interface {
	GetAll(req GetScansRequest) ([]types.Scan, error)
	Get(scanID string) (types.Scan, error)
	Create(scan types.Scan) (types.Scan, error)
	Update(scanID string, scan types.Scan) (types.Scan, error)
	Delete(scanID string) error
}

type GetScanFindingsRequest struct {
	ScanID *string `json:"scanID"`
}

type DeleteScanFindingsRequest struct {
	ScanID *string `json:"scanID"`
}

type ScanFindingStore interface {
	GetAll(req GetScanFindingsRequest) ([]types.ScanFinding, error)
	Get(findingID string) (types.ScanFinding, error)
	CreateMany(scanID string, findings ...types.ScanFinding) ([]types.ScanFinding, error)

	// Update is not needed since we only keep data in-memory for analytical purposes
	// Update(findingID string, finding ScanFinding) (ScanFinding, error)

	Delete(req DeleteScanFindingsRequest) error
}

type Store interface {
	Scans() ScanStore
	ScanFindings() ScanFindingStore
}
