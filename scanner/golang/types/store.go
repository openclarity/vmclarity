package types

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type PreconditionFailedError struct {
	Reason string
}

func (e *PreconditionFailedError) Error() string {
	return "Precondition failed: " + e.Reason
}

type GetScansRequest struct {
	State *string `json:"state,omitempty"`
}

type ScanStore interface {
	GetAll(req GetScansRequest) ([]Scan, error)
	Get(scanID string) (Scan, error)
	Create(scan Scan) (Scan, error)
	Update(scanID string, scan Scan) (Scan, error)
	Delete(scanID string) error
}

type GetScanFindingsRequest struct {
	ScanID string `json:"scanID"`
}

type ScanFindingStore interface {
	GetAll(req GetScanFindingsRequest) ([]ScanFinding, error)
	Get(findingID string) (ScanFinding, error)
	Create(finding ScanFinding) (ScanFinding, error)

	// Update is not needed since we only keep data in-memory for analytical purposes
	// Update(findingID string, finding ScanFinding) (ScanFinding, error)

	Delete(findingID string) error
}

type Store interface {
	Scans() ScanStore
	ScanFindings() ScanFindingStore
}
